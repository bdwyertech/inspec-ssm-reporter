package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
	// AWS SDK
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ssm"
)

type ExecutionSummary struct{}
type InSpecReport map[string]interface{}

// Determine Severity
func getSeverity(severity float64) string {
	switch {
	case severity >= 0.0 && severity < 0.4:
		return "LOW"
	case severity >= 0.4 && severity < 0.7:
		return "HIGH"
	case severity >= 0.7:
		return "CRITICAL"
	default:
		return "CRITICAL"
	}
}

func main() {
	instance_id, ok := os.LookupEnv("AWS_SSM_INSTANCE_ID")
	if ok !=true {
		log.Fatal("Unable to find environment variable AWS_SSM_INSTANCE_ID: make sure this script is executed by SSM Agent")
	}

	region, ok := os.LookupEnv("AWS_SSM_REGION_NAME")
	if ok !=true {
		log.Fatal("Unable to find environment variable AWS_SSM_REGION_NAME: make sure this script is executed by SSM Agent")
	}

	// Read the JSON from STDIN
    var report InSpecReport

    err := json.NewDecoder(os.Stdin).Decode(&report)
    if err != nil {
        log.Fatal(err)
    }

	// Derive ExecutionID from the PWD
	// PWD is something like: /var/lib/amazon/ssm/INSTANCE_ID/document/orchestration/EXECUTION_ID/downloads/
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	execution_id := filepath.Base(filepath.Dir(pwd))

    // Construct the Compliance Items
    var items []*ssm.ComplianceItemEntry

	compliant := 0;
	non_compliant := 0;
	compliant_by_sev := make(map[string]int)
	non_compliant_by_sev := make(map[string]int)

	for k, v := range report {
		if k == "profiles" {
			profiles := v.([]interface{})
			for _, profile := range profiles {
				profile := profile.(map[string]interface{})
				controls := profile["controls"].([]interface{})
				for _, control := range controls {
					control := control.(map[string]interface{})
					results := control["results"].([]interface{})
					for _, result := range results {
						result := result.(map[string]interface{})
						severity := getSeverity(control["impact"].(float64))
						status := ""
						if result["status"] == "passed" {
							status = "COMPLIANT"
							compliant++
							compliant_by_sev[severity]++
						} else if result["status"] == "failed" {
							status = "NON_COMPLIANT"
							non_compliant++
							non_compliant_by_sev[severity]++
						} else {
							continue
						}

						items = append(items, &ssm.ComplianceItemEntry{
							Id: aws.String(fmt.Sprintf("%s-%d", control["id"], len(items))),
							Severity: aws.String(severity),
							Status: aws.String(status),
							Title: aws.String(fmt.Sprintf("%s : %s", control["title"], result["code_desc"])),
						})
					}
				}
			}
		}
	}
	// DEBUG
	fmt.Println(items)

    // Create a SSM client with additional configuration
	svc := ssm.New(session.New(), aws.NewConfig().WithRegion(region))

	// Submit the Compliance Report
	svc.PutComplianceItems(&ssm.PutComplianceItemsInput{
		ResourceId: 	&instance_id,
	  	ResourceType: 	aws.String("ManagedInstance"),
  		ComplianceType: aws.String("Custom:InSpec"),
  		ExecutionSummary: &ssm.ComplianceExecutionSummary{
  			ExecutionId: aws.String(execution_id),
  			ExecutionTime: aws.Time(time.Now()),
  			ExecutionType: aws.String("Command"),
  		},
  		Items: items,
	})

	fmt.Printf("Completed InSpec checks and put %d compliant and %d non-compliant items", compliant, non_compliant)
}
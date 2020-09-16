package main

import (
	"encoding/json"
	"flag"
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

func InSpecToCompliance(report InSpecReport) (items []*ssm.ComplianceItemEntry) {
	compliant := 0
	non_compliant := 0
	compliant_by_sev := make(map[string]int)
	non_compliant_by_sev := make(map[string]int)
	for _, profile := range report.Profiles {
		for _, control := range profile.Controls {
			for _, result := range control.Results {
				severity := getSeverity(control.Impact)
				status := ""
				if result.Status == "passed" {
					status = "COMPLIANT"
					compliant++
					compliant_by_sev[severity]++
				} else if result.Status == "failed" {
					status = "NON_COMPLIANT"
					non_compliant++
					non_compliant_by_sev[severity]++
				} else {
					continue
				}

				items = append(items, &ssm.ComplianceItemEntry{
					Id:       aws.String(fmt.Sprintf("%s-%d", control.ID, len(items))),
					Severity: aws.String(severity),
					Status:   aws.String(status),
					Title:    aws.String(fmt.Sprintf("%s : %s", control.Title, result.CodeDesc)),
				})
			}
		}
	}
	// DEBUG
	// fmt.Printf("%d compliant and %d non-compliant items", compliant, non_compliant)
	// fmt.Println(items)
	return
}

func main() {
	flag.Parse()

	if versionFlag {
		showVersion()
		os.Exit(0)
	}

	instance_id, ok := os.LookupEnv("AWS_SSM_INSTANCE_ID")
	if !ok {
		log.Fatal("Unable to find environment variable AWS_SSM_INSTANCE_ID: make sure this is executed by SSM Agent")
	}

	region, ok := os.LookupEnv("AWS_SSM_REGION_NAME")
	if !ok {
		log.Fatal("Unable to find environment variable AWS_SSM_REGION_NAME: make sure this is executed by SSM Agent")
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
	items := InSpecToCompliance(report)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            *aws.NewConfig().WithRegion(region).WithCredentialsChainVerboseErrors(true),
		SharedConfigState: session.SharedConfigDisable,
	}))

	// Create a SSM client with additional configuration
	svc := ssm.New(sess)

	// Submit the Compliance Report
	svc.PutComplianceItems(&ssm.PutComplianceItemsInput{
		ResourceId:     &instance_id,
		ResourceType:   aws.String("ManagedInstance"),
		ComplianceType: aws.String("Custom:InSpec"),
		ExecutionSummary: &ssm.ComplianceExecutionSummary{
			ExecutionId:   aws.String(execution_id),
			ExecutionTime: aws.Time(time.Now()),
			ExecutionType: aws.String("Command"),
		},
		Items: items,
	})

	fmt.Println("Completed InSpec checks")
}

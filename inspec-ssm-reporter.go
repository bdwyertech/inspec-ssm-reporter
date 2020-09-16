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
	items := report.ToComplianceItems()

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

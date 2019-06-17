# inspec-ssm-reporter

[![Build Status](https://travis-ci.org/bdwyertech/inspec-ssm-reporter.svg)](https://travis-ci.org/bdwyertech/inspec-ssm-reporter)

###  Overview
This is a utility to transform InSpec JSON into an AWS Compliance document

* Parses JSON from STDIN
* Transforms into an AWS Compliance Document
* Reports to SSM Compliance

### Background
The default AWS-provided pattern leverages the Ruby environment provided by ChefDK and installs `aws-sdk-ssm` directly from Rubygems.  Installing ChefDK for this is heavy-handed and not ideal at scale.  Additionally, the scripts pull installation packages directly from the Internet which does not work in an air-gapped environment.

The goal here is to deploy InSpec by itself (much smaller package) and leverage this static Golang binary to handle the compliance reporting.  This removes the need for a Ruby environment

#### AWS Equivalent
* http://aws-ssm-us-east-1.s3.amazonaws.com/inspec/report_compliance
##### Calling Scripts
* http://aws-ssm-us-east-1.s3.amazonaws.com/inspec/run_inspec.ps1
* http://aws-ssm-us-east-1.s3.amazonaws.com/inspec/run_inspec.sh

### Usage
#### Linux
```bash
inspec exec . --reporter json | inspec-ssm-reporter
if [ $? -ne 0 ]; then
  echo "Failed to execute InSpec tests: see stderr"
  EXITCODE=2
fi
```

#### Windows
```powershell
$results=inspec exec . --reporter json 2> errors.txt
$results | inspec-ssm-reporter
if(!$?) {
  Write-Host "Failed to execute InSpec tests: see stderr"
  $EXITCODE=2
}
```

### Development
1. Use gvm under WSL
2. gvm install go1.12.6

### InSpec JSON Model
* https://github.com/inspec/inspec/blob/master/test/unit/mock/reporters/json_output

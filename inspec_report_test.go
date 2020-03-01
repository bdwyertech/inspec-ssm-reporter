// Encoding: UTF-8

package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestInspecStructParse(t *testing.T) {
	url := "https://raw.githubusercontent.com/inspec/inspec/master/test/fixtures/reporters/json_output"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Could not download test fixture %s", url)
	}
	defer resp.Body.Close()

	report := InSpecReport{}

	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		t.Fatal("Could not unmarshal InSpec report")
	}

	for _, profile := range report.Profiles {
		for _, control := range profile.Controls {
			for _, result := range control.Results {
				severity := getSeverity(control.Impact)
				status := ""
				if result.Status == "passed" {
					status = "COMPLIANT"
				} else if result.Status == "failed" {
					status = "NON_COMPLIANT"
				} else {
					continue
				}
				t.Logf("%v - %v - %v", status, severity, result)
			}
		}
	}
	// t.Fatal("")
}

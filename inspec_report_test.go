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

	items := InSpecToCompliance(report)
	t.Log(items)
}

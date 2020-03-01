// Encoding: UTF-8

package main

// https://github.com/inspec/inspec/blob/master/lib/inspec/reporters/json.rb

type InSpecReport struct {
	Platform struct {
		Name     string `json:"name,omitempty"`
		Release  string `json:"release,omitempty"`
		TargetId string `json:"target_id,omitempty"`
	} `json:"platform"`
	Profiles []struct {
		Name           string `json:"name"`
		Version        string `json:"version"`
		Sha256         string `json:"sha256"`
		Title          string `json:"title"`
		Maintainer     string `json:"maintainer"`
		Summary        string `json:"summary"`
		License        string `json:"license"`
		Copyright      string `json:"copyright"`
		CopyrightEmail string `json:"copyright_email"`

		Supports []struct {
			OsFamily string `json:"os-family,omitempty"`
			OsName   string `json:"os-name,omitempty"`
			Release  string `json:"release,omitempty"`
		} `json:"supports"`

		Attributes []interface{} `json:"attributes"`

		Depends []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"depends"`

		Groups []struct {
			ID       string   `json:"id"`
			Controls []string `json:"controls"`
			Title    string   `json:"title,omitempty"`
		} `json:"groups"`

		Controls []struct {
			ID           string `json:"id"`
			Title        string `json:"title"`
			Desc         string `json:"desc"`
			Descriptions []struct {
				Label string `json:"label"`
				Data  string `json:"data"`
			} `json:"descriptions"`
			Impact         float64       `json:"impact"`
			Refs           []interface{} `json:"refs"`
			Tags           struct{}      `json:"tags,omitempty"`
			Code           string        `json:"code"`
			SourceLocation struct {
				Line int    `json:"line"`
				Ref  string `json:"ref"`
			} `json:"source_location"`
			WaiverData struct {
			} `json:"waiver_data"`

			Results []struct {
				Status      string  `json:"status"`
				CodeDesc    string  `json:"code_desc,omitempty"`
				RunTime     float64 `json:"run_time,omitempty"`
				StartTime   string  `json:"start_time,omitempty"`
				Resource    string  `json:"resource,omitempty"`
				SkipMessage string  `json:"skip_message,omitempty"`
				Message     string  `json:"message,omitempty"`
				Exception   string  `json:"exception,omitempty"`
				Backtrace   string  `json:"backtrace,omitempty"`
			} `json:"results"`
		} `json:"controls"`
	} `json:"profiles"`

	Statistics struct {
		Duration float64 `json:"duration"`
	} `json:"statistics"`

	Version string `json:"version"`
}

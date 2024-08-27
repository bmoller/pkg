package aur

import (
	"encoding/json"
)

type responseBase struct {
	Version     int    `json:"version"`
	Type        string `json:"type"`
	ResultCount int    `json:"resultcount"`
}

type errorResponse struct {
	responseBase
	Error   string        `json:"error"`
	Results []interface{} `json:"results,omitempty"`
}

type searchResponse struct {
	responseBase
	Results []packageBasic `json:"results"`
}

type infoResponse struct {
	responseBase
	Results []packageDetailed `json:"results"`
}

type packageBasic struct {
	ID             int
	Name           string
	Description    string
	PackageBaseID  int
	PackageBase    string
	Maintainer     string
	NumVotes       int
	Popularity     json.Number
	FirstSubmitted int
	LastModified   int
	OutOfDate      string
	Version        string
	URLPath        string
	URL            string
}

type response struct {
	ResultCount int               `json:"resultcount"`
	Type        string            `json:"type"`
	Version     int               `json:"version"`
	Error       string            `json:"error"`
	Results     []packageDetailed `json:"results"`
}

func (r *response) String() string {
	if j, err := json.Marshal(r); err != nil {
		return ""
	} else {
		return string(j)
	}
}

type packageDetailed struct {
	packageBasic
	Submitter     string
	License       []string
	Depends       []string
	MakeDepends   []string
	OptDepends    []string
	CheckDepends  []string
	Provides      []string
	Conflicts     []string
	Replaces      []string
	Groups        []string
	Keywords      []string
	CoMaintainers []string
}

func (p packageDetailed) String() string {
	if j, err := json.Marshal(p); err != nil {
		return ""
	} else {
		return string(j)
	}
}

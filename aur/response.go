package aur

import (
	"encoding/json"
)

type result struct {
	ResultCount int          `json:"resultcount"`
	Type        string       `json:"type"`
	Version     int          `json:"version"`
	Error       string       `json:"error"`
	Results     []aurPackage `json:"results"`
}

/*
String provides the JSON string representation of an AUR result by marshalling
the object using Go's json package. If the function returns an error the
returned string is empty.
*/
func (r *result) String() string {
	if j, err := json.Marshal(r); err != nil {
		return ""
	} else {
		return string(j)
	}
}

type aurPackage struct {
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
	Submitter      string
	License        []string
	Depends        []string
	MakeDepends    []string
	OptDepends     []string
	CheckDepends   []string
	Provides       []string
	Conflicts      []string
	Replaces       []string
	Groups         []string
	Keywords       []string
	CoMaintainers  []string
}

/*
String provides the JSON string representation of an AUR result by marshalling
the object using Go's json package. If the function returns an error the
returned string is empty.
*/
func (p aurPackage) String() string {
	if j, err := json.Marshal(p); err != nil {
		return ""
	} else {
		return string(j)
	}
}

package aur

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// fallback value used if terminal width lookup fails
const defaultWidth = 80

// number of spaces to indent for line continuations
const leftPadding = 18

// time layout used by pacman
const timeLayout = "Mon 02 Jan 2006 03:04:05 PM MST"

/*
A result represents a response from the AUR API.
The fields are tagged for support of marshalling using Go's json package.
*/
type result struct {
	ResultCount int          `json:"resultcount"` // number of results in Results
	Type        string       `json:"type"`        // AUR response type: error, info, or search
	Version     int          `json:"version"`     // server-side version of the API
	Error       string       `json:"error"`       // error message returned by the API, if any
	Results     []aurPackage `json:"results"`     // slice of packages returned for a search or info request
}

/*
String provides the JSON string representation of an AUR result by marshalling the object using Go's json package.
If the function returns an error the returned string is empty.
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
	FirstSubmitted int // Unix timestamp
	LastModified   int // Unix timestamp
	OutOfDate      int // Unix timestamp
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
String provides the JSON string representation of an AUR result by marshalling the object using Go's json package.
If the function returns an error the returned string is empty.
*/
func (p aurPackage) String() string {
	if j, err := json.Marshal(p); err != nil {
		return ""
	} else {
		return string(j)
	}
}

/*
Formatted builds a string representation of p in a human-readable format.
It handles terminals of varying width and should provide output regardless of the current terminal size.
The overall formatting is very similar to, but not an exact match for,
the format used by pacman to display information about a package.
*/
func (p aurPackage) Formatted() string {
	s := "Name            : " + p.Name + "\n"
	s += "Version         : " + p.Version + "\n"
	s += "Description     : " + p.Description + "\n"
	s += "URL             : " + p.URL + "\n"
	s += "Licenses        : "
	s += printArray(p.License)
	s += "Groups          : "
	s += printArray(p.Groups)
	s += "Provides        : "
	s += printArray(p.Provides)
	s += "Depends On      : "
	s += printArray(p.Depends)
	s += "Optional Deps   : "
	if len(p.OptDepends) == 0 {
		s += "None\n"
	} else {
		s += strings.Join(p.OptDepends, "\n                  ") + "\n"
	}
	s += "Conflicts With  : "
	s += printArray(p.Conflicts)
	s += "Replaces        : "
	s += printArray(p.Replaces)
	s += "Maintainer      : " + p.Maintainer + "\n"
	s += "Last Modified   : "
	s += time.Unix(int64(p.LastModified), 0).Format(timeLayout) + "\n"

	return s
}

/*
printArray builds a string that will display well-formatted in the current terminal.
It attempts to look up the current terminal's width, but if this fails a default fallback value is used.
*/
func printArray(a []string) string {
	if len(a) == 0 {
		return "None\n"
	}

	// build a prefix of the appropriate number of spaces
	padding := ""
	for range leftPadding {
		padding += " "
	}

	// get the width of the terminal
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// set a safe fallback value if the lookup fails
		width = defaultWidth
	}

	// build a string with line breaks where needed
	output := ""
	length := leftPadding
	for _, s := range a {
		if length+len(s)+2 > width {
			output += "\n" + padding
			length = leftPadding
		}
		output += s + "  "
		length += len(s) + 2
	}
	output += "\n"

	return output
}

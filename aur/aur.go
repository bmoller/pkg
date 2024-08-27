package aur

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// AURHost is the HTTP scheme and domain of the AUR.
const aurHost = "https://aur.archlinux.org"

// AURInfoPath is the URL path of the AUR information endpoint.
const aurInfoPath = "/rpc/v5/info"

// AURSearchPath is the URL path of the AUR search endpoint.
const aurSearchPath = "/rpc/v5/search"

// A SearchType is the kind of AUR search to perform. It determines which fields of
// packages a search term will match against.
type SearchType int

const (
	NameDesc     SearchType = iota // match name or description
	Name                           // match package names only
	Maintainer                     // match package maintainers
	Depends                        // match package dependencies
	MakeDepends                    // match dependencies required to build a package
	OptDepends                     // match optional dependencies of a package
	CheckDepends                   // match dependencies required to check a package

	Fake
)

var queryKeys = map[SearchType]string{
	NameDesc:     "name-desc",
	Name:         "name",
	Maintainer:   "maintainer",
	Depends:      "depends",
	MakeDepends:  "makedepends",
	OptDepends:   "optdepends",
	CheckDepends: "checkdepends",
	Fake:         "fake",
}

func Search(keyword string, by SearchType) (result *response, err error) {
	target, err := url.Parse(aurHost)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %e", err)
	}
	target = target.JoinPath(aurSearchPath, keyword)
	query := target.Query()
	query.Set("by", queryKeys[by])
	target.RawQuery = query.Encode()

	resp, err := http.Get(target.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make search request: %e", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %e", err)
	}
	result = new(response)
	if err = json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response JSON: %e", err)
	}

	return
}

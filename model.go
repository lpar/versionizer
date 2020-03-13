package versionize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"
)

// MustParseURL parses a URL, or panics if it can't be parsed
func MustParseURL(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(fmt.Errorf("can't parse URL: %w", err))
	}
	return u
}

// URI is a wrapper to enable JSON unmarshaling of Go url.URL values
type URI struct {
	url.URL
}

// UnmarshalJSON implements the Unmarshaler interface.
func (u *URI) UnmarshalJSON(b []byte) error {
	return u.URL.UnmarshalBinary(bytes.Trim(b, `"`))
}

// Metadata defines the metadata about a piece of software.
type Metadata struct {
	Title string
	Description string
	Authors string
	Vendor string
	// Licenses is expected to be an SPDX license expression, see https://spdx.org/spdx-specification-21-web-version#h.jxpfx0ykyb60
	Licenses string
	Version string
	Revision string
	Created time.Time
	Information *URI
	Documentation *URI
	SourceCode *URI
}

// Load reads a JSON manifest into a Metadata struct
func Load(jsonfile string) (Metadata, error) {
	m := Metadata{}
	dat, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(dat, &m)
	return m, err
}
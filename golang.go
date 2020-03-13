package versionize

import (
	"bufio"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

// This file contains functions used to write metadata as Go code.

const goTemplate = `
package {{ .Package }}

const {{ .Prefix -}} Title = "{{ escape .M.Title }}"
const {{ .Prefix -}} Description = "{{ escape .M.Description }}"
const {{ .Prefix -}} Authors = "{{ escape .M.Authors }}"
const {{ .Prefix -}} Vendor = "{{ escape .M.Vendor }}"
const {{ .Prefix -}} Licenses = "{{ escape .M.Licenses }}"
const {{ .Prefix -}} Version = "{{ escape .M.Version }}"
const {{ .Prefix -}} Revision = "{{ escape .M.Revision }}"
const {{ .Prefix -}} Created = "{{ rfc3339 .M.Created }}"
const {{ .Prefix -}} Information = "{{ .M.Information }}"
const {{ .Prefix -}} Documentation = "{{ .M.Documentation }}"
const {{ .Prefix -}} SourceCode = "{{ .M.SourceCode }}"
`

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"rfc3339": func (t time.Time) string {
			return t.Format(time.RFC3339)
		},
		"escape": func (s string) string {
			return strings.ReplaceAll(s, `"`, `\"`)
		},
	}
}

// WriteGoCode returns the metadata as compilable Go code.
// The code is written with the package name specified, and each constant is prefixed with the specified prefix.
func (m Metadata) WriteGoCode(out io.Writer, packname string, prefix string) error {
	data := struct{
		Package string
		Prefix string
		M Metadata
	}{
		Package: packname,
		Prefix: prefix,
		M: m,
	}
	t, err := template.New("version.go").Funcs(getFuncMap()).Parse(goTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(out, data)
	return err
}

// WriteGoFile writes the metadata to compilable Go source code, with the specified filename, Go package declaration
// and prefix to add to each constant.
func (m Metadata) WriteGoFile(fname string, packname string, prefix string) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = m.WriteGoCode(w, packname, prefix)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}
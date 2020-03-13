package versionize

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// This file contains functions used to write Open Container Initiative metadata into a Containerfile.

const ociTemplate = `LABEL org.opencontainers.image.title="{{ escape .Title }}"
LABEL org.opencontainers.image.description="{{ escape .Description }}"
LABEL org.opencontainers.image.authors="{{ escape .Authors }}"
LABEL org.opencontainers.image.vendor="{{ escape .Vendor }}"
LABEL org.opencontainers.image.licenses="{{ escape .Licenses }}"
LABEL org.opencontainers.image.version="{{ escape .Version }}"
LABEL org.opencontainers.image.revision="{{ escape .Revision }}"
LABEL org.opencontainers.image.created="{{ rfc3339 .Created }}"
LABEL org.opencontainers.image.information="{{ .Information }}"
LABEL org.opencontainers.image.documentation="{{ .Documentation }}"
LABEL org.opencontainers.image.sourceCode="{{ .SourceCode }}"
`

func (m Metadata) WriteOCI(w *bufio.Writer) error {
	t, err := template.New("version.go").Funcs(getFuncMap()).Parse(ociTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(w, m)
	return err
}

func (m Metadata) SpliceInto(fname string, start string, end string) error {
	var startMatched bool
  var  endMatched bool
	fin, err := os.Open(fname)
	if err != nil {
		return err
	}
	fdir := filepath.Dir(fname)
	tmpfile, err := ioutil.TempFile(fdir, fname)
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	w := bufio.NewWriter(tmpfile)
	s := bufio.NewScanner(fin)
	for s.Scan() {
		line := s.Text()
		if _, err = w.WriteString(line); err != nil {
			return err
		}
		if _, err = w.WriteRune('\n'); err != nil {
			return err
		}
		if strings.Contains(line, start) {
			startMatched = true
			// Write replacement
			err = m.WriteOCI(w)
			if err != nil {
				return err
			}
			// Skip until end
			for s.Scan() {
				lin := s.Text()
				if strings.Contains(lin, end) {
					endMatched = true
					if _, err = w.WriteString(lin); err != nil {
						return err
					}
					if _, err = w.WriteRune('\n'); err != nil {
						return err
					}
					break
				}
			}
		}
	}
	if !startMatched {
		return fmt.Errorf("didn't find start substring '%s' in file %s", start, fname)
	}
	if !endMatched {
		return fmt.Errorf("didn't find end substring '%s' in file %s", end, fname)
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	err = tmpfile.Close()
	if err != nil {
		return err
	}
  // Now replace the original
	err = os.Remove(fname)
	if err != nil {
		return err
	}
	err = os.Rename(tmpfile.Name(), fname)
	return err
}

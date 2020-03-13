
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lpar/versionize"
)

var manifestFile = flag.String("manifest", "Manifest.json", "JSON metadata manifest to read")
var goFile = flag.String("go", "", "Go file to write")
var goPackage = flag.String("package", "main", "Go package name to write")
var goPrefix = flag.String("prefix", "meta", "Prefix for Go constant names")
var ociFile = flag.String("oci", "", "Containerfile/Dockerfile to update")
var ociStart = flag.String("start", "# Begin metadata", "String marking start of metadata in Containerfile")
var ociEnd = flag.String("end", "# End metadata", "String marking end of metadata in Containerfile")

func main() {
	flag.Parse()
	var m versionize.Metadata
	if *manifestFile != "" {
		m = loadMetadata(*manifestFile)
	}
	if *goFile != "" {
		writeGoCode(*goFile, m)
	}
	if *ociFile != "" {
		writeOCI(*ociFile, m)
	}
}

func writeOCI(arg string, m versionize.Metadata) {
	err := m.SpliceInto(arg, *ociStart, *ociEnd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't update container file: %v\n", err)
		os.Exit(1)
	}
}

func writeGoCode(arg string, m versionize.Metadata) {
	err := m.WriteGoFile(arg, *goPackage, *goPrefix)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't load metadata: %v\n", err)
		os.Exit(2)
	}
}

func loadMetadata(arg string) versionize.Metadata {
	m, err := versionize.Load(arg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't load metadata: %v\n", err)
		os.Exit(3)
	}
	m.Revision, err = versionize.GitRevision()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(4)
	}
	m.Version, err = versionize.GitVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(5)
	}
	m.Created = time.Now()
	return m
}

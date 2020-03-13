
# versionize

This is a utility for reading metadata from your project's tree using Git, and from a JSON file. That information can
then be written out into a compilable `.go` file of constants, and/or spliced in to a Containerfile or Dockerfile.

You can specify the package to generate the Go constants in, and a prefix for their names.

In a container file, the data is written following the [OCI image specification][oci] standard 
[predefined annotation keys][keys], and you can specify the substrings to look for when updating the file.

See the `example/` directory for an example manifest file and container file. For an example of the Go output, 
see `cmd/versionize/metadata.go`.

See `magefile.go` for an example of using this program with [Mage][mage].

[oci]: https://github.com/opencontainers/image-spec
[keys]: https://github.com/opencontainers/image-spec/blob/master/annotations.md#pre-defined-annotation-keys
[mage]: https://magefile.org/

## Caveats

 - No unit test coverage yet, it's an afternoon hack written for my own use.
 - Only works with Git.
 - Assumes you're doing something other than just `go build` to build your application.


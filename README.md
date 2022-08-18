# schnittstelle

schnittstelle *(german for "interface")* is a CLI tool to generate an interface from the methods
implemented for a struct in a source tree.

## Installation

Simply compile and install the tool on your system using `go install`.
```
GOBIN=/usr/local/bin go install github.com/zekrotja/schnittstelle/v3/cmd/schnittstelle@latest
```

Alternatively, you can simply download latest binaries from the 
[releases page](https://github.com/zekroTJA/schnittstelle/releases) or artifacts from 
[latest CI builds](https://github.com/zekroTJA/schnittstelle/actions/workflows/releases.yml).

## Usage

```
Usage of schnittstelle:
  -interface string
        The name of the result interface (name of struct when not specified)
  -out string
        Output file (if not specified, output will be piped to Stdout)
  -package string
        Package name ingested in output
  -pool uint
        Number of files which can be searched through simultaneously (default 10)
  -root string
        Root directory (default ".")
  -struct string
        Name of the struct
```
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
Usage: schnittstelle --struct STRUCT [--root ROOT] [--interface INTERFACE] [--package PACKAGE] [--out OUT] [--inject INJECT] [--import IMPORT] [--pool POOL]

Options:
  --struct STRUCT, -s STRUCT
                         Name of the struct
  --root ROOT, -r ROOT   Root directory [default: .]
  --interface INTERFACE, -i INTERFACE
                         The name of the result interface (name of struct when not specified)
  --package PACKAGE, -p PACKAGE
                         Package name ingested in output
  --out OUT, -o OUT      Output file (if not specified, output will be piped to Stdout)
  --inject INJECT        Inject code lines into the output code.
  --import IMPORT        Add import lines to the output.
  --pool POOL            Number of files which can be searched through simultaneously [default: 5]
  --help, -h             display this help and exit
```
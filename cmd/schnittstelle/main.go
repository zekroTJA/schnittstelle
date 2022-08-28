package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/zekrotja/schnittstelle"
)

type Args struct {
	Struct    string   `arg:"-s,--struct,required" help:"Name of the struct"`
	Root      string   `arg:"-r,--root" help:"Root directory" default:"."`
	Interface string   `arg:"-i,--interface" help:"The name of the result interface (name of struct when not specified)"`
	Package   string   `arg:"-p,--package" help:"Package name ingested in output"`
	Out       string   `arg:"-o,--out" help:"Output file (if not specified, output will be piped to Stdout)"`
	Inject    []string `arg:"--inject,separate" help:"Inject code lines into the output code."`
	Imports   []string `arg:"--import,separate" help:"Add import lines to the output."`
	PoolSize  uint     `arg:"--pool" help:"Number of files which can be searched through simultaneously" default:"5"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	interfaceName := args.Interface
	if interfaceName == "" {
		interfaceName = args.Struct
	}

	signatures, err := schnittstelle.Extract(
		args.Struct, args.Root, int(args.PoolSize))
	if err != nil {
		fmt.Println("Error: Extracting signatures:", err.Error())
		return
	}

	var output io.WriteCloser = os.Stdout
	outFile := args.Out
	if outFile != "" {
		if !strings.HasSuffix(outFile, ".go") {
			outFile += ".go"
		}
		pathTo := filepath.Dir(outFile)
		_, err := os.Stat(pathTo)
		if os.IsNotExist(err) {
			err = os.MkdirAll(pathTo, 0770)
			if err != nil {
				fmt.Println("Error: Creating path to output file:", err)
				return
			}
		} else if err != nil {
			fmt.Println("Error: Checking output path:", err)
			return
		}
		output, err = os.Create(outFile)
		if err != nil {
			fmt.Println("Error: Opening output file:", err)
			return
		}
		defer output.Close()
	}

	inject := args.Inject
	if len(args.Imports) != 0 {
		inject = make([]string, 0, len(args.Inject)+len(args.Imports)+2)
		inject = append(inject, "import (")
		for _, imp := range args.Imports {
			inject = append(inject, fmt.Sprintf("\t%s",
				formatImport(imp)))
		}
		inject = append(inject, ")")
		inject = append(inject, args.Inject...)
	}

	err = schnittstelle.Assemble(
		interfaceName, args.Package, inject,
		signatures, output)
	if err != nil {
		fmt.Println("Error: Assembling output:", err.Error())
		return
	}
}

func formatImport(v string) string {
	split := strings.SplitN(v, " ", 2)

	var prefix, pack string
	if len(split) == 2 {
		prefix = split[0]
		pack = split[1]
	} else {
		pack = split[0]
	}

	if pack[0] != '"' && pack[len(pack)-1] != '"' {
		pack = `"` + pack + `"`
	}

	return prefix + " " + pack
}

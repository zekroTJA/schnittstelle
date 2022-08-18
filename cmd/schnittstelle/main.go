package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zekrotja/schnittstelle"
)

var (
	fStructName = flag.String("struct", "", "Name of the struct")
	fRoot       = flag.String("root", ".", "Root directory")
	fInterface  = flag.String("interface", "",
		"The name of the result interface (name of struct when not specified)")
	fPackage = flag.String("package", "", "Package name ingested in output")
	fOut     = flag.String("out", "",
		"Output file (if not specified, output will be piped to Stdout)")
	fPoolSize = flag.Uint("pool", 10,
		"Number of files which can be searched through simultaneously")
)

func main() {
	flag.Parse()

	structName := *fStructName

	if structName == "" {
		fmt.Println("Error: struct name must be given")
		return
	}

	interfaceName := *fInterface
	if interfaceName == "" {
		interfaceName = structName
	}

	signatures, err := schnittstelle.Extract(
		structName, *fRoot, int(*fPoolSize))
	if err != nil {
		fmt.Println("Error: Extracting signatures:", err.Error())
		return
	}

	var output io.WriteCloser = os.Stdout
	outFile := *fOut
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

	err = schnittstelle.Assemble(interfaceName, *fPackage, signatures, output)
	if err != nil {
		fmt.Println("Error: Assembling output:", err.Error())
		return
	}
}

package schnittstelle

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
)

var (
	ErrPathIsFile = errors.New("root path is not a directory")
)

// Extract walks through all files recursively starting from
// codePath creating an index of Go source files which are
// no unit test files.
//
// Then, it looks through all files line-by-line to find
// methods which have the given structName as receiver. This
// applies either to both receivers by value or by reference.
// Only methods which are exported are included in the result
// set.
//
// The poolSize defines the number of files which can be opened
// and searched through simultaneously.
//
// The resulting method signatures are then returned as an
// array of strings. The result set is sorted alphabetically.
func Extract(structName string, codePath string, poolSize int) ([]string, error) {
	index := make([]string, 0, 100)

	stat, err := os.Stat(codePath)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, ErrPathIsFile
	}

	err = filepath.WalkDir(codePath,
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(d.Name()) != ".go" {
				return nil
			}
			if strings.HasSuffix(d.Name(), "_test.go") {
				return nil
			}

			index = append(index, path)
			return nil
		})
	if err != nil {
		return nil, err
	}

	if len(index) == 0 {
		return nil, errors.New("The specified directory does not contain any Go files")
	}

	pool := NewWorkerpool(poolSize)

	cErr := make(chan error)
	signatures := make([][]string, 0, len(index))
	errors := MultiError{}

	var mtx sync.Mutex

	go func() {
		for {
			select {
			case res := <-pool.Results():
				if sig, ok := res.([]string); ok {
					mtx.Lock()
					signatures = append(signatures, sig)
					mtx.Unlock()
				}
			case err := <-cErr:
				mtx.Lock()
				errors = append(errors, err)
				mtx.Unlock()
			}
		}
	}()

	for _, filePath := range index {
		pool.Push(func(workerId int, params ...interface{}) interface{} {
			sig, err := FindMethodsInFile(params[0].(string), structName)
			if err != nil {
				cErr <- err
				return nil
			}
			return sig
		}, filePath)
	}

	pool.Close()
	pool.WaitBlocking()
	time.Sleep(10 * time.Millisecond)

	var size int
	for _, sigs := range signatures {
		size += len(sigs)
	}

	signaturesFlat := make(sort.StringSlice, 0, size)
	for _, sigs := range signatures {
		signaturesFlat = append(signaturesFlat, sigs...)
	}

	sort.Sort(signaturesFlat)

	return signaturesFlat, errors.ToError()
}

// FindMethodsInFile reads the given file line by line
// to look for methods having the given structName as
// receiver and returns the list of results in an array.
func FindMethodsInFile(filePath string, structName string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	structNamePointer := "*" + structName

	signatures := []string{}

	var methodSig string
	var comment bool
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())

		if comment {
			if !strings.HasSuffix(txt, "*/") {
				continue
			}
			comment = false
		}
		if strings.HasPrefix(txt, "//") {
			continue
		}
		if strings.HasPrefix(txt, "/*") {
			comment = true
			continue
		}

		var sig string
		if methodSig == "" {
			if !strings.HasPrefix(txt, "func (") {
				continue
			}
			recEnd := strings.IndexRune(txt, ')')
			receiver := txt[len("func ("):recEnd]
			spacePos := strings.IndexRune(receiver, ' ')
			if spacePos != -1 {
				receiver = receiver[spacePos+1:]
			}
			if receiver != structName && receiver != structNamePointer {
				continue
			}

			sig = txt[recEnd+2:]
			if !unicode.IsUpper(rune(sig[0])) {
				continue
			}
		} else {
			sig = txt
		}

		last := sig[len(sig)-1]
		if last == '}' {
			sigEnd := strings.LastIndex(sig, "{")
			sig = sig[:sigEnd+1]
		} else if last != '{' {
			methodSig += sig
			if last != '(' {
				methodSig += " "
			}
			continue
		}

		methodSig += sig[:len(sig)-2]
		methodSig = strings.ReplaceAll(methodSig, ", )", ")")

		signatures = append(signatures, methodSig)

		methodSig = ""
	}

	return signatures, nil
}

// Assemble takes an array of method signatures and builds
// an interface with the defined interfaceName and writes it
// to w. If a packageName is specified, the package header
// with the given package name is added to the output. When
// inject is specified, the contets will be added to the
// output in between the package statement and the interface
// injection.
func Assemble(
	interfaceName string,
	packageName string,
	inject []string,
	signatures []string,
	w io.Writer,
) (err error) {
	if packageName != "" {
		_, err = fmt.Fprintf(w, "package %s\n\n", packageName)
		if err != nil {
			return err
		}
	}

	if len(inject) != 0 {
		for _, i := range inject {
			i = strings.ReplaceAll(i, "\\n", "\n")
			i = strings.ReplaceAll(i, "\\t", "\t")
			fmt.Fprintf(w, "%s\n", i)
		}
		fmt.Fprint(w, "\n")
	}

	_, err = fmt.Fprintf(w, "type %s interface {\n", interfaceName)
	if err != nil {
		return err
	}

	for _, sig := range signatures {
		_, err = fmt.Fprintf(w, "\t%s\n", sig)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprint(w, "}\n")
	return err
}

package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

// Generate creates CodeResult for given Sysl definitions as Proto message (pb.Module)
func Generate(module *pb.Module, pkg string) (CodeResult, error) {
	apps := module.GetApps()
	if len(apps) != 1 {
		return CodeResult{}, fmt.Errorf("cannot handle more than 1 application")
	}
	buffer := bytes.NewBufferString(fmt.Sprintf("package %s\n\n", pkg))
	for _, app := range apps {
		interf, err := GenInterface(app)
		if err != nil {
			return CodeResult{}, err
		}
		buffer.WriteString(interf)
		types, err := GenTypes(app)
		if err != nil {
			return CodeResult{}, err
		}
		buffer.WriteString(types)
	}
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return CodeResult{}, err
	}
	return CodeResult{Interface: string(b)}, nil
}

// CodeResult contains REST and interface golang source files' contents as strings
type CodeResult struct {
	Rest      string
	Interface string
}

// GetPackage extracts package name from output directory
func GetPackage(outDir string) string {
	l := strings.Split(outDir, "/")
	return l[len(l)-1]
}

// LineName contains name, such as type or endpoint name and corresponding line
type LineName struct {
	name string
	line int32
}

// SortLineNames sorts a slice of LineName in place and returns a slice of
// sorted names
func SortLineNames(lineNames []LineName) []string {
	sort.Slice(lineNames, func(i, j int) bool {
		return lineNames[i].line < lineNames[j].line
	})
	size := len(lineNames)
	result := make([]string, size)
	for i := 0; i < size; i++ {
		result[i] = lineNames[i].name
	}
	return result
}

package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

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

// Generate creates CodeResult for given Sysl definitions as Proto message (pb.Module)
func Generate(module *pb.Module, pkg string) (CodeResult, error) {
	apps := module.GetApps()
	if len(apps) != 1 {
		return CodeResult{}, fmt.Errorf("cannot handle more than 1 application")
	}
	buffer := bytes.NewBufferString(fmt.Sprintf("package %s\n\n", pkg))
	for _, app := range apps {
		jsonSep := ""
		if attr, ok := app.Attrs["json_property_separator"]; ok {
			jsonSep = attr.GetS()
		}
		types, err := GenTypes(app.GetTypes(), jsonSep)
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

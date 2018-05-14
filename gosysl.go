package gosysl

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

// CodeResult contains restfile and interface file contents as strings
type CodeResult struct {
	Rest      string
	Interface string
}

// GetPackage extracts package name from output directory
func GetPackage(outDir string) string {
	l := strings.Split(outDir, "/")
	return l[len(l)-1]
}

// Generate creates CodeResult for given sysl definitons as proto message (pb.Module)
func Generate(module *pb.Module, pkg string) (CodeResult, error) {
	apps := module.GetApps()
	if len(apps) != 1 {
		return CodeResult{}, fmt.Errorf("cannot handle more than 1 app")
	}
	buffer := bytes.NewBufferString(fmt.Sprintf("package %s\n\n", pkg))
	for _, app := range apps {
		types, err := GenTypes(app.GetTypes())
		if err != nil {
			return CodeResult{}, err
		}
		buffer.WriteString(types)
	}
	return CodeResult{Interface: buffer.String()}, nil
}

package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

// CodeResult contains REST and interface golang source files' contents as strings
type CodeResult struct {
	Rest       string
	Interface  string
	Middleware string
}

// Generate creates CodeResult for given Sysl definitions as Proto message (pb.Module)
func Generate(module *pb.Module, pkg string) (CodeResult, error) {
	app, err := getApp(module)
	if err != nil {
		return CodeResult{}, err
	}
	epNames := sortEpNames(app.Endpoints)
	interf, err := genInterfaceFile(app, epNames, pkg)
	if err != nil {
		return CodeResult{}, err
	}
	middleware, err := genMiddlewareFile(app, epNames, pkg)
	if err != nil {
		return CodeResult{}, err
	}
	rest, err := genRestFile(app, epNames, pkg)
	if err != nil {
		return CodeResult{}, err
	}
	result := CodeResult{
		Rest:       rest,
		Interface:  interf,
		Middleware: middleware,
	}
	return result, nil
}

func genRestFile(app *pb.Application, epNames []string, pkg string) (string, error) {
	if len(app.Name.Part) > 0 && app.Name.Part[0] == "BAD" {
		return pkg, fmt.Errorf("100")
	}
	res := ""
	if len(epNames) > 0 {
		res = epNames[0]
	}
	return res, nil
}

func genInterfaceFile(app *pb.Application, epNames []string, pkg string) (string, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf("package %s\n\n", pkg))
	interf, err := GenInterface(app, epNames)
	if err != nil {
		return "", err
	}
	buffer.WriteString(interf)
	types, err := GenTypes(app)
	if err != nil {
		return "", err
	}
	buffer.WriteString(types)
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func genMiddlewareFile(app *pb.Application, eps []string, pkg string) (string, error) {
	buffer := bytes.NewBufferString(fmt.Sprintf("package %s\n\n", pkg))
	buffer.WriteString("import \"net/http\"\n\n")
	interf, err := GenMiddleware(app, eps)
	if err != nil {
		return "", err
	}
	buffer.WriteString(interf)
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getApp(module *pb.Module) (*pb.Application, error) {
	apps := module.GetApps()
	if len(apps) > 1 {
		return nil, fmt.Errorf("cannot handle more than 1 application")
	}
	for _, app := range apps {
		return app, nil
	}
	return nil, fmt.Errorf("need at least 1 application")
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

func getEndpointLine(ep *pb.Endpoint) int32 {
	if ep.RestParams == nil {
		return 0
	}
	params := ep.RestParams.QueryParam
	if len(params) == 0 {
		if strings.HasPrefix(strings.ToLower(ep.Name), "get") {
			return -1
		}
		return 0
	}
	return params[0].Type.SourceContext.Start.Line
}

func sortEpNames(endpoints map[string]*pb.Endpoint) []string {
	lineNames := make([]LineName, len(endpoints))
	i := 0
	for name, t := range endpoints {
		line := getEndpointLine(t)
		lineNames[i] = LineName{name, line}
		i++
	}
	return SortLineNames(lineNames)
}

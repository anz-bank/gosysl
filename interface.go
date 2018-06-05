package gosysl

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

var reMethodRemove = regexp.MustCompile(`[{}\s]`)
var reMethodSeparate = regexp.MustCompile(`[._,#-]`)

// GetMethodName creates the interface method name from pattern path
func GetMethodName(ep *pb.Endpoint) string {
	if n, ok := ep.Attrs["method_name"]; ok {
		return n.GetS()
	}
	name := reMethodRemove.ReplaceAllLiteralString(ep.Name, "")
	name = reMethodSeparate.ReplaceAllLiteralString(name, "/")

	fields := strings.Split(name, "/")
	for i, field := range fields {
		fields[i] = strings.Title(strings.ToLower(field))
	}
	return strings.Join(fields, "")
}

var curlyRe = regexp.MustCompile(`^\s*{\s*(\w+)\s*<:\s*(\w+)\s*}\s*$`)

func getParams(ep *pb.Endpoint) (string, error) {
	if ep.RestParams == nil {
		return "", nil
	}
	params := make([]string, 0, 8)
	for _, param := range ep.RestParams.QueryParam {
		name := param.Name
		typeStr, _, err := GetType(param.Type)
		if err != nil {
			return "", err
		}
		matches := curlyRe.FindAllStringSubmatch(typeStr, -1)
		if len(matches) == 1 {
			name = matches[0][1]
			typeStr = matches[0][2]
			params = append(params, fmt.Sprintf("%s %s", name, typeStr))
		} else {
			params = append([]string{fmt.Sprintf("%s %s", name, typeStr)}, params...)
		}
	}
	for _, param := range ep.Param {
		typeStr := param.Type.GetTypeRef().Ref.Appname.Part[0]
		params = append(params, fmt.Sprintf("%s %s", param.Name, typeStr))
	}
	return strings.Join(params, ", "), nil
}

func getReturnTypes(ep *pb.Endpoint) (string, error) {
	for _, s := range ep.Stmt {
		if s.GetAction() != nil && s.GetAction().GetAction() == "return" {
			// simple return type without value
			return "error", nil
		}
		if s.GetRet() != nil {
			retType := s.GetRet().GetPayload()
			return fmt.Sprintf("(%s, error)", retType), nil
		}
	}
	return "", fmt.Errorf("return missing in endpoint %s", ep.String())
}

func writeMethod(w io.Writer, ep *pb.Endpoint) error {
	if attr, ok := ep.Attrs["method_doc"]; ok {
		fmt.Fprintf(w, "\n// %s \n", attr.GetS())
	}
	name := GetMethodName(ep)
	params, err := getParams(ep)
	if err != nil {
		return err
	}
	returnTypes, err := getReturnTypes(ep)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s(%s) %s\n", name, params, returnTypes)
	return nil
}

// WriteInterface creates for methods called in REST endpoints.
func WriteInterface(w io.Writer, app *pb.Application, epNames []string) error {
	if attr, ok := app.Attrs["interface_doc"]; ok {
		fmt.Fprintf(w, "// %s \n", attr.GetS())
	}
	interfaceName := "Storer"
	if attr, ok := app.Attrs["interface"]; ok {
		interfaceName = strings.Title(attr.GetS())
	}
	fmt.Fprintf(w, "type %s interface {\n", interfaceName)
	for _, name := range epNames {
		if err := writeMethod(w, app.Endpoints[name]); err != nil {
			return err
		}
	}
	fmt.Fprintln(w, "}")
	return nil
}

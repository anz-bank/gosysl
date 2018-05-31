package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

var reMethodRemove = regexp.MustCompile(`[{}\s]`)
var reMethodSeparate = regexp.MustCompile(`[._,#-]`)

// GenMethodName creates the interface method name from pattern path
func GenMethodName(ep *pb.Endpoint) string {
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

func genParams(ep *pb.Endpoint) (string, error) {
	if ep.RestParams == nil {
		return "", nil
	}
	params := make([]string, 0, 8)
	for _, param := range ep.RestParams.QueryParam {
		name := param.Name
		typeStr, _, err := GenType(param.Type)
		if err != nil {
			return "", err
		}
		matches := curlyRe.FindAllStringSubmatch(typeStr, -1)
		if len(matches) == 1 {
			name = matches[0][1]
			typeStr = matches[0][2]
		}
		params = append(params, fmt.Sprintf("%s %s", name, typeStr))
	}
	for _, param := range ep.Param {
		typeStr := param.Type.GetTypeRef().Ref.Appname.Part[0]
		params = append(params, fmt.Sprintf("%s %s", param.Name, typeStr))
	}
	return strings.Join(params, ", "), nil
}

func genReturnTypes(ep *pb.Endpoint) (string, error) {
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

func genMethod(ep *pb.Endpoint) (string, error) {
	var buffer bytes.Buffer
	if attr, ok := ep.Attrs["method_doc"]; ok {
		buffer.WriteString(fmt.Sprintf("\n// %s \n", attr.GetS()))
	}
	name := GenMethodName(ep)
	params, err := genParams(ep)
	if err != nil {
		return "", err
	}
	returnTypes, err := genReturnTypes(ep)
	if err != nil {
		return "", err
	}
	buffer.WriteString(fmt.Sprintf("%s(%s) %s\n", name, params, returnTypes))
	return buffer.String(), nil
}

// GenInterface creates for methods called in REST endpoints.
func GenInterface(app *pb.Application, epNames []string) (string, error) {
	var buffer bytes.Buffer
	if attr, ok := app.Attrs["interface_doc"]; ok {
		buffer.WriteString(fmt.Sprintf("// %s \n", attr.GetS()))
	}
	interfaceName := "Storer"
	if attr, ok := app.Attrs["interface"]; ok {
		interfaceName = strings.Title(attr.GetS())
	}
	buffer.WriteString(fmt.Sprintf("type %s interface {\n", interfaceName))

	for _, name := range epNames {
		method, err := genMethod(app.Endpoints[name])
		if err != nil {
			return "", err
		}
		buffer.WriteString(method)
	}

	buffer.WriteString("}\n")
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

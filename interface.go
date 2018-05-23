package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"regexp"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

func getEndpointLine(ep *pb.Endpoint) int32 {
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

func genMethodName(ep *pb.Endpoint) string {
	if n, ok := ep.Attrs["method_name"]; ok {
		return n.GetS()
	}
	name := strings.Replace(ep.Name, " ", "", -1)
	name = strings.Replace(name, "{", "", -1)
	name = strings.Replace(name, "}", "", -1)

	name = strings.Replace(name, "-", "/", -1)
	name = strings.Replace(name, "_", "/", -1)
	name = strings.Replace(name, ".", "/", -1)
	name = strings.Replace(name, ",", "/", -1)
	name = strings.Replace(name, "#", "/", -1)

	fields := strings.Split(name, "/")
	for i, field := range fields {
		fields[i] = strings.Title(strings.ToLower(field))
	}
	return strings.Join(fields, "")
}

var curlyRe = regexp.MustCompile(`^\s*{\s*(\w+)\s*<:\s*(\w+)\s*}\s*$`)

func genParams(ep *pb.Endpoint) (string, error) {
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
	name := genMethodName(ep)
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
func GenInterface(app *pb.Application) (string, error) {
	var buffer bytes.Buffer
	if attr, ok := app.Attrs["interface_doc"]; ok {
		buffer.WriteString(fmt.Sprintf("// %s \n", attr.GetS()))
	}
	interfaceName := "Storer"
	if attr, ok := app.Attrs["interface"]; ok {
		interfaceName = strings.Title(attr.GetS())
	}
	buffer.WriteString(fmt.Sprintf("type %s interface {\n", interfaceName))

	names := sortEpNames(app.Endpoints)
	for _, name := range names {
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

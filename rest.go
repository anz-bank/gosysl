package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

// GenMiddleware creates interface returning required middleware functions
// for REST endpoints
func GenMiddleware(app *pb.Application, epNames []string) (string, error) {
	middlewares := make([]string, 0, len(epNames))
	middlewareSet := make(map[string]interface{}, len(epNames))
	for _, name := range epNames {
		if m, ok := app.Endpoints[name].Attrs["middleware"]; ok {
			if _, ok2 := middlewareSet[m.GetS()]; !ok2 {
				middlewareSet[m.GetS()] = struct{}{}
				middlewares = append(middlewares, m.GetS())
			}
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString("type Middleware interface {\n")
	for _, m := range middlewares {
		s := fmt.Sprintf("%s () []func(next http.Handler) http.Handler\n", m)
		buffer.WriteString(s)
	}
	buffer.WriteString("}\n")
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// GenRest creates the contextkeys, routes and handlers for actual REST handlers
func GenRest(app *pb.Application, epNames []string) (string, error) {
	buffer := bytes.NewBufferString(genContextKeys(app, epNames))
	return buffer.String(), nil
}

func genContextKeys(app *pb.Application, epNames []string) string {
	keys := getContextKeys(app, epNames)
	if len(keys) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("// Keys for Context lookup\nconst (\n")
	s := fmt.Sprintf("%s ContextKeyType = iota\n", getContextKey(keys[0]))
	buffer.WriteString(s)
	for i := 1; i < len(keys); i++ {
		buffer.WriteString(getContextKey(keys[i]) + "\n")
	}
	buffer.WriteString(")\n")
	return buffer.String()
}

func getContextKey(param string) string {
	return strings.Title(param) + "Key"
}

func getContextKeys(app *pb.Application, epNames []string) []string {
	set := make(map[string]interface{}, len(epNames))
	result := make([]string, 0, len(epNames))
	for _, name := range epNames {
		rp := app.Endpoints[name].RestParams
		if rp == nil {
			continue
		}
		for _, qp := range rp.QueryParam {
			if _, ok := set[qp.Name]; !ok {
				set[qp.Name] = struct{}{}
				result = append(result, qp.Name)
			}
		}
	}
	return result
}

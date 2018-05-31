package gosysl

import (
	"bytes"
	"fmt"
	"go/format"

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

	buffer := bytes.NewBufferString("type Middleware interface {\n")
	for _, m := range middlewares {
		s := fmt.Sprintf("%s () []func(next http.Handler) http.Handler\n", m)
		buffer.WriteString(s) // nolint: gas
	}
	buffer.WriteString("}\n") // nolint: gas
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

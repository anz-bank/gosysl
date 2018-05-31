package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

type route struct {
	methods         map[string]string
	middleware      string
	keys            []string
	queryParams     []string
	postPayloadType string
	putPayloadType  string
}

type routes struct {
	paths   []string
	content map[string]*route
}

var validHTTPMethods = map[string]interface{}{
	"Get":    struct{}{},
	"Put":    struct{}{},
	"Post":   struct{}{},
	"Delete": struct{}{},
}

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
	var buffer bytes.Buffer
	buffer.WriteString(genContextKeys(app, epNames))
	r, err := getRoutes(app, epNames)
	if err != nil {
		return "", err
	}
	buffer.WriteString(genNewRestHandler(r))
	buffer.WriteString(genHandlers(r))
	b, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func genHandlers(rs routes) string {
	var buffer bytes.Buffer
	for _, path := range rs.paths {
		r := rs.content[path]
		if handler, ok := r.methods["Get"]; ok {
			buffer.WriteString(genGet(handler, r))
		}
		if handler, ok := r.methods["Post"]; ok {
			buffer.WriteString(genPost(handler, r))
		}
		if handler, ok := r.methods["Put"]; ok {
			buffer.WriteString(genPut(handler, r))
		}
		if handler, ok := r.methods["Delete"]; ok {
			buffer.WriteString(genDelete(handler, r))
		}
	}
	return buffer.String()
}

func getQueryVar(queryParam string) string {
	return "query" + strings.Title(queryParam)
}

func getHandlerHead(handler string, keys []string, queryParams []string) string {
	var buffer bytes.Buffer
	s := "func (rh *RestHandler) handle%s(w http.ResponseWriter, r *http.Request) {\n"
	buffer.WriteString(fmt.Sprintf(s, handler))
	for _, key := range keys {
		s = "	%s := r.Context().Value(%s).(string)\n"
		buffer.WriteString(fmt.Sprintf(s, key, getContextKey(key)))
	}
	for _, qp := range queryParams {
		s = "	%s := r.URL.Query().Get(\"%s\")\n"
		buffer.WriteString(fmt.Sprintf(s, getQueryVar(qp), qp))
	}
	return buffer.String()
}

const errBoiler = `if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}`

func genGet(handler string, r *route) string {
	var buffer bytes.Buffer
	buffer.WriteString(getHandlerHead(handler, r.keys, r.queryParams))
	params := strings.Join(append(r.keys, r.queryParams...), ", ")
	s := `	result, err := rh.storer.%s(%s)
	%s
	render.JSON(w, r, result)
}` + "\n\n"
	buffer.WriteString(fmt.Sprintf(s, handler, params, errBoiler))
	return buffer.String()
}

func genDelete(handler string, r *route) string {
	var buffer bytes.Buffer
	buffer.WriteString(getHandlerHead(handler, r.keys, r.queryParams))
	params := strings.Join(append(r.keys, r.queryParams...), ", ")
	s := `	if err := rh.storer.%s(%s); err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.NoContent(w, r)
}` + "\n\n"
	buffer.WriteString(fmt.Sprintf(s, handler, params))
	return buffer.String()
}

const payloadBoiler = `	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}` + "\n"

func genPutPost(handler string, r *route, post bool) string {
	var buffer bytes.Buffer
	buffer.WriteString(getHandlerHead(handler, r.keys, r.queryParams))
	p := append(r.keys, r.queryParams...)
	p = append(p, "payload")
	params := strings.Join(p, ", ")
	payloadType := r.putPayloadType
	status := ""
	if post {
		payloadType = r.postPayloadType
		status = "\nrender.Status(r, http.StatusCreated)"
	}

	s := `	var payload %s
	%s
	result, err := rh.storer.%s(%s)
	%s %s
	render.JSON(w, r, result)
}` + "\n\n"
	str := fmt.Sprintf(s,
		payloadType,
		payloadBoiler,
		handler,
		params,
		errBoiler,
		status)
	buffer.WriteString(str)
	return buffer.String()
}

func genPut(handler string, r *route) string {
	return genPutPost(handler, r, false)
}

func genPost(handler string, r *route) string {
	return genPutPost(handler, r, true)
}

func genNewRestHandler(r routes) string {
	var buffer bytes.Buffer
	buffer.WriteString(`// NewRestHandler creates a new Handler persisting data to Storer.
func NewRestHandler(s Storer, m Middleware) RestHandler {
	r := chi.NewRouter()
	r.Use(m.Root()...)
	rh := RestHandler{s, r}` + "\n\n")
	buffer.WriteString(genRoutes(r))
	buffer.WriteString(`return rh
}` + "\n\n")
	return buffer.String()
}

func genRoutes(r routes) string {
	var buffer bytes.Buffer
	for _, path := range r.paths {
		str := fmt.Sprintf("r.Route(\"%s\", func(r chi.Router) {\n", path)
		buffer.WriteString(str)
		for _, key := range r.content[path].keys {
			ctxKey := getContextKey(key)
			str = fmt.Sprintf("r.Use(makeContextSaver(%s, \"%s\"))\n", ctxKey, key)
			buffer.WriteString(str)
		}
		middleware := r.content[path].middleware
		if middleware != "" {
			buffer.WriteString(fmt.Sprintf("r.Use(m.%s()...)\n", middleware))
		}
		methods := r.content[path].methods
		for _, m := range []string{"Get", "Post", "Put", "Delete"} {
			handler, ok := methods[m]
			if ok {
				buffer.WriteString(fmt.Sprintf("r.%s(\"/\", rh.handle%s)\n", m, handler))
			}
		}
		buffer.WriteString("})\n")
	}
	return buffer.String()
}

func getRoutes(app *pb.Application, epNames []string) (routes, error) {
	paths := make([]string, 0, len(epNames)/2)
	content := make(map[string]*route, len(epNames)/2)
	for _, name := range epNames {
		fields := strings.Split(name, " ")
		if len(fields) != 2 {
			msg := `expect "GET|POST|etc path/path" as endpoint name (%s) `
			return routes{}, fmt.Errorf(msg, name)
		}
		endpoint := app.Endpoints[name]
		httpPath := fields[1]
		httpMethod := strings.Title(strings.ToLower(fields[0]))
		if _, ok := validHTTPMethods[httpMethod]; !ok {
			return routes{}, fmt.Errorf("invalid HTTP Method (%s)", httpMethod)
		}
		if _, ok := content[httpPath]; !ok {
			middleware := ""
			if m, ok := endpoint.Attrs["middleware"]; ok {
				middleware = m.GetS()
			}
			content[httpPath] = &route{
				methods:     map[string]string{},
				middleware:  middleware,
				keys:        getPatternParams(endpoint),
				queryParams: getQueryParams(endpoint),
			}
			paths = append(paths, httpPath)
		}
		interfaceMethod := GenMethodName(endpoint)
		content[httpPath].methods[httpMethod] = interfaceMethod
		if httpMethod == "Put" {
			t := getPayloadType(endpoint)
			content[httpPath].putPayloadType = t
		}
		if httpMethod == "Post" {
			content[httpPath].postPayloadType = getPayloadType(endpoint)
		}
	}
	return routes{paths, content}, nil
}

func getPayloadType(ep *pb.Endpoint) string {
	if len(ep.Param) != 1 || ep.Param[0].Type.GetTypeRef() == nil {
		return ""
	}
	return ep.Param[0].Type.GetTypeRef().Ref.Appname.Part[0]
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

func getPatternParams(ep *pb.Endpoint) []string {
	return getParams(ep, false)
}

func getQueryParams(ep *pb.Endpoint) []string {
	return getParams(ep, true)
}

func getParams(ep *pb.Endpoint, queryType bool) []string {
	result := make([]string, 0, len(ep.RestParams.QueryParam))
	if ep.RestParams == nil {
		return result
	}
	for _, qp := range ep.RestParams.QueryParam {
		if (queryType && qp.Type.GetTypeRef() != nil) ||
			(!queryType && qp.Type.GetTypeRef() == nil) {
			result = append(result, qp.Name)
		}
	}
	return result

}

func getContextKeys(app *pb.Application, epNames []string) []string {
	set := make(map[string]interface{}, len(epNames))
	result := make([]string, 0, len(epNames))
	for _, name := range epNames {
		for _, qp := range getPatternParams(app.Endpoints[name]) {
			if _, ok := set[qp]; !ok {
				set[qp] = struct{}{}
				result = append(result, qp)
			}
		}
	}
	return result
}

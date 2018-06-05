package gosysl

import (
	"fmt"
	"io"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

type route struct {
	methods         map[string]string
	middleware      string
	keys            []string
	queryParams     map[string][]string
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

// WriteMiddleware writes interface returning required middleware functions
// for REST endpoints
func WriteMiddleware(w io.Writer, app *pb.Application, epNames []string) {
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

	fmt.Fprintln(w, `// Middleware holds the middleware accessor methods for the REST API`)
	fmt.Fprintln(w, `type Middleware interface {`)
	for _, m := range middlewares {
		fmt.Fprintf(w, "%s() []func(next http.Handler) http.Handler\n", m)
	}
	fmt.Fprintln(w, "Root() []func(next http.Handler) http.Handler")
	fmt.Fprintln(w, `}`)
}

// WriteRest creates the contextkeys, routes and handlers for actual REST handlers
func WriteRest(w io.Writer, app *pb.Application, epNames []string) error {
	writeContextKeys(w, app, epNames)
	r, err := getRoutes(app, epNames)
	if err != nil {
		return err
	}
	writeNewRestHandler(w, r)
	writeHandlers(w, r)
	return nil
}

func writeHandlers(w io.Writer, rs routes) {
	for _, path := range rs.paths {
		r := rs.content[path]
		if handler, ok := r.methods["Get"]; ok {
			writeGet(w, handler, r)
		}
		if handler, ok := r.methods["Post"]; ok {
			writePost(w, handler, r)
		}
		if handler, ok := r.methods["Put"]; ok {
			writePut(w, handler, r)
		}
		if handler, ok := r.methods["Delete"]; ok {
			writeDelete(w, handler, r)
		}
	}
}

func writeHandlerHead(w io.Writer, handler string, keys []string, queryParams []string) {
	format := "func (rh *RestHandler) handle%s(w http.ResponseWriter, r *http.Request) {\n"
	fmt.Fprintf(w, format, handler)
	for _, key := range keys {
		format = "	%s := r.Context().Value(%s).(string)\n"
		fmt.Fprintf(w, format, key, getContextKey(key))
	}
	for _, qp := range queryParams {
		format = "	%s := r.URL.Query().Get(\"%s\")\n"
		fmt.Fprintf(w, format, qp, qp)
	}
}

const errBoiler = `if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}`

func writeGet(w io.Writer, handler string, r *route) {
	writeHandlerHead(w, handler, r.keys, r.queryParams["Get"])
	params := strings.Join(append(r.keys, r.queryParams["Get"]...), ", ")
	s := `	result, err := rh.storer.%s(%s)
	%s
	render.JSON(w, r, result)
}` + "\n\n"
	fmt.Fprintf(w, s, handler, params, errBoiler)
}

func writeDelete(w io.Writer, handler string, r *route) {
	writeHandlerHead(w, handler, r.keys, r.queryParams["Delete"])
	params := strings.Join(append(r.keys, r.queryParams["Delete"]...), ", ")
	s := `	if err := rh.storer.%s(%s); err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.NoContent(w, r)
}` + "\n\n"
	fmt.Fprintf(w, s, handler, params)
}

const payloadBoiler = `	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}` + "\n"

func writePut(w io.Writer, handler string, r *route) {
	writeHandlerHead(w, handler, r.keys, r.queryParams["Put"])
	p := append(r.keys, r.queryParams["Put"]...)
	p = append(p, "payload")
	params := strings.Join(p, ", ")
	s := `	var payload %s
	%s
	result, err := rh.storer.%s(%s)
	%s
	render.JSON(w, r, result)
}` + "\n\n"
	fmt.Fprintf(w, s, r.putPayloadType, payloadBoiler, handler, params, errBoiler)
}

func writePost(w io.Writer, handler string, r *route) {
	writeHandlerHead(w, handler, r.keys, r.queryParams["Post"])
	p := append(r.keys, r.queryParams["Post"]...)
	p = append(p, "payload")
	params := strings.Join(p, ", ")

	s := `	var payload %s
	%s
	result, err := rh.storer.%s(%s)
	%s
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, result)
}` + "\n\n"
	fmt.Fprintf(w, s, r.postPayloadType, payloadBoiler, handler, params, errBoiler)
}

func writeNewRestHandler(w io.Writer, r routes) {
	fmt.Fprint(w, `// NewRestHandler creates a new Handler persisting data to Storer.
func NewRestHandler(s Storer, m Middleware) RestHandler {
	r := chi.NewRouter()
	r.Use(m.Root()...)
	rh := RestHandler{s, r}`+"\n\n")
	writeRoutes(w, r)
	fmt.Fprint(w, "return rh \n} \n\n")
}

func writeRoutes(w io.Writer, r routes) {
	for _, path := range r.paths {
		fmt.Fprintf(w, `r.Route("%s", func(r chi.Router) {`+"\n", path)
		for _, key := range r.content[path].keys {
			ctxKey := getContextKey(key)
			fmt.Fprintf(w, "r.Use(makeContextSaver(%s, \"%s\"))\n", ctxKey, key)
		}
		middleware := r.content[path].middleware
		if middleware != "" {
			fmt.Fprintf(w, "r.Use(m.%s()...)\n", middleware)
		}
		methods := r.content[path].methods
		for _, m := range []string{"Get", "Post", "Put", "Delete"} {
			handler, ok := methods[m]
			if ok {
				fmt.Fprintf(w, "r.%s(\"/\", rh.handle%s)\n", m, handler)
			}
		}
		fmt.Fprint(w, "})\n")
	}
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
		httpMethod := strings.Title(strings.ToLower(fields[0]))
		if _, ok := validHTTPMethods[httpMethod]; !ok {
			return routes{}, fmt.Errorf("invalid HTTP Method (%s)", httpMethod)
		}
		endpoint := app.Endpoints[name]
		httpPath := fields[1]
		if _, ok := content[httpPath]; !ok {
			middleware := ""
			if m, ok := endpoint.Attrs["middleware"]; ok {
				middleware = m.GetS()
			}
			content[httpPath] = &route{
				methods:     map[string]string{},
				middleware:  middleware,
				keys:        getPatternParams(endpoint),
				queryParams: make(map[string][]string, 4),
			}
			paths = append(paths, httpPath)
		}
		interfaceMethod := GetMethodName(endpoint)
		content[httpPath].methods[httpMethod] = interfaceMethod
		content[httpPath].queryParams[httpMethod] = getQueryParams(endpoint)
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

func writeContextKeys(w io.Writer, app *pb.Application, epNames []string) {
	keys := getContextKeys(app, epNames)
	if len(keys) == 0 {
		return
	}
	fmt.Fprintln(w, `// Keys for Context lookup`)
	fmt.Fprintln(w, `const (`)
	fmt.Fprintf(w, "%s ContextKeyType = iota\n", getContextKey(keys[0]))
	for i := 1; i < len(keys); i++ {
		fmt.Fprintln(w, getContextKey(keys[i]))
	}
	fmt.Fprintln(w, ")")
}

func getContextKey(param string) string {
	return strings.Title(param) + "Key"
}

func getPatternParams(ep *pb.Endpoint) []string {
	if ep == nil || ep.RestParams == nil {
		return nil
	}
	result := make([]string, 0, len(ep.RestParams.QueryParam))
	for _, qp := range ep.RestParams.QueryParam {
		if qp.Type.GetTypeRef() == nil {
			result = append([]string{qp.Name}, result...)
		}
	}
	return result
}

func getQueryParams(ep *pb.Endpoint) []string {
	if ep == nil || ep.RestParams == nil {
		return nil
	}
	result := make([]string, 0, len(ep.RestParams.QueryParam))
	for _, qp := range ep.RestParams.QueryParam {
		if qp.Type.GetTypeRef() != nil {
			result = append(result, qp.Name)
		}
	}
	return result
}

func getContextKeys(app *pb.Application, epNames []string) []string {
	set := make(map[string]struct{}, len(epNames))
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

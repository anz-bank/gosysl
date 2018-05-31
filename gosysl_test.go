package gosysl

import (
	"io/ioutil"
	"testing"

	"github.com/anz-bank/gosysl/pb"
	"github.com/golang/protobuf/proto"
	testifyAssert "github.com/stretchr/testify/assert"
)

func TestEnd2End(tt *testing.T) {
	assert := testifyAssert.New(tt)
	data, err := ioutil.ReadFile("example/example.pb")
	assert.NoError(err)
	module := &pb.Module{}
	err = proto.Unmarshal(data, module)
	assert.NoError(err)
	result, err := Generate(module, "mypkg")
	assert.NoError(err)
	assert.Equal(expectedInterface, result.Interface)
	assert.Equal(expectedMiddleware, result.Middleware)
	assert.Equal(expectedRest, result.Rest)

	// failing gofmt
	_, err = Generate(module, "BAD PACKAGE NAME")
	assert.Error(err)
}

func TestGetPackage(tt *testing.T) {
	assert := testifyAssert.New(tt)

	assert.Equal("x", GetPackage("x"))
	assert.Equal("y", GetPackage("x/y"))
}

func TestGetApp(tt *testing.T) {
	assert := testifyAssert.New(tt)

	module := &pb.Module{}
	module.Apps = map[string]*pb.Application{}
	module.Apps["1"] = nil
	module.Apps["2"] = nil
	_, err := getApp(module)
	assert.Error(err)
}

var expectedInterface = `package mypkg

// Storer abstracts all required RefData persistence and retrieval
type Storer interface {

	// DataSet
	GetKeys() (Keys, error)
	CreateDataSet(ds DataSetPayload) (Key, error)
	GetDataSetName(key string, queryTime string) (KeyName, error)
	PutDataSetName(key string, np NamePayload) (KeyName, error)

	// Data
	GetData(key string, queryTime string) (Data, error)
	PutData(key string, dp DataPayload) (Data, error)
	GetDataWithStart(startTime string, key string) (Data, error)
	PutDataWithStart(startTime string, key string, dp DataPayload) (Data, error)
	DeleteData(startTime string, key string) error

	// Schema
	GetSchema(key string, queryTime string) (Schema, error)
	PutSchema(key string, sp SchemaPayload) (Schema, error)
	GetSchemaWithStart(startTime string, key string) (Schema, error)
	PutSchemaWithStart(startTime string, key string, sp SchemaPayload) (Schema, error)
	DeleteSchema(startTime string, key string) error

	// Admin
	DeleteDataSet(key string) error
	GetStartTimes(key string) (Times, error)
	GetCreationTimes(key string) (CreationTimes, error)
	GetRestriction(key string) (Restriction, error)
	PutRestriction(key string, r Restriction) (Restriction, error)
	PutSubscription(key string, s Subscription) (Subscription, error)
	DeleteSubscription(key string, s Subscription) error
}

// Data holds JSON data valid from StartTime created at CreationTime
type Data struct {
	StartTime    string ` + "`json:\"start-time\"`" + `
	JSONData     string ` + "`json:\"data\"`" + `
	CreationTime string ` + "`json:\"creation-time\"`" + `
}

// Schema holds JSON Schema to validate Data against, a name for key creation.
type Schema struct {
	StartTime    string ` + "`json:\"start-time\"`" + `
	JSONSchema   string ` + "`json:\"schema\"`" + `
	CreationTime string ` + "`json:\"creation-time\"`" + `
}

// Subscription holds external endpoint values for change notification
type Subscription struct {
	URL          string ` + "`json:\"url\"`" + `
	SecreteToken string ` + "`json:\"secrete-token\"`" + `
}

// Restriction contains scope access restriction and frozen times for schema and data.
type Restriction struct {
	SchemaFrozenUntil string   ` + "`json:\"schema-frozen-until\"`" + `
	DataFrozenUntil   string   ` + "`json:\"data-frozen-until\"`" + `
	ReadScopes        []string ` + "`json:\"read-scopes\"`" + `
	ReadWriteScopes   []string ` + "`json:\"read-write-scopes\"`" + `
	AdminScopes       []string ` + "`json:\"admin-scopes\"`" + `
}

// Keys is JSON result type for getKeys in REST API
type Keys struct {
	Keys []string ` + "`json:\"keys\"`" + `
}

// Key is JSON result type for createDataSet in REST API
type Key struct {
	Key string ` + "`json:\"key\"`" + `
}

// KeyName is JSON result type for get and put dataDetNamre in REST API
type KeyName struct {
	Key  string ` + "`json:\"key\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

// Times contains schema and data times, used to get StartTimes for both
type Times struct {
	Data   []string ` + "`json:\"data-times\"`" + `
	Schema []string ` + "`json:\"schema\"`" + `
}

// CreationStartTime contains start and creation time for a schema or data snapshot
type CreationStartTime struct {
	CreationTime string ` + "`json:\"creation-time\"`" + `
	StartTime    string ` + "`json:\"start-time\"`" + `
}

// CreationTimes contains schema and data times maps, used to StartTime to CreationTims
type CreationTimes struct {
	Data   map[string]CreationStartTime ` + "`json:\"data-time-map\"`" + `
	Schema map[string]CreationStartTime ` + "`json:\"schema\"`" + `
}

// DataSetPayload is JSON payload on REST API request to create new data set
type DataSetPayload struct {
	Name         string      ` + "`json:\"name\"`" + `
	StartTimeStr string      ` + "`json:\"start-time\"`" + `
	JSONSchema   interface{} ` + "`json:\"schema\"`" + `
}

// NamePayload is JSON payload on REST API request to update data set name
type NamePayload struct {
	Name string ` + "`json:\"name\"`" + `
}

// DataPayload is JSON payload on REST API request to update data
type DataPayload struct {
	Data interface{} ` + "`json:\"data\"`" + `
}

// SchemaPayload is JSON payload on REST API request to update schema
type SchemaPayload struct {
	Schema interface{} ` + "`json:\"schema\"`" + `
}

// UpdateEvent holds all information necessary to post to subscribes
type UpdateEvent struct {
	Key       string      ` + "`json:\"key\"`" + `
	StartTime string      ` + "`json:\"start-time\"`" + `
	Data      interface{} ` + "`json:\"data\"`" + `
	Schema    interface{} ` + "`json:\"schema\"`" + `
	Deleted   bool        ` + "`json:\"deleted\"`" + `
}
`

var expectedMiddleware = `package mypkg

import "net/http"

type Middleware interface {
	AuthorizeRoot() []func(next http.Handler) http.Handler
	AuthorizeDataSet() []func(next http.Handler) http.Handler
	AuthorizeAdmin() []func(next http.Handler) http.Handler
}
`

var expectedRest = "package mypkg\n\n" + restPrefix + `
// Keys for Context lookup
const (
	KeyKey ContextKeyType = iota
	StartTimeKey
)

// NewRestHandler creates a new Handler persisting data to Storer.
func NewRestHandler(s Storer, m Middleware) RestHandler {
	r := chi.NewRouter()
	r.Use(m.Root()...)
	rh := RestHandler{s, r}

	r.Route("/api", func(r chi.Router) {
		r.Use(m.AuthorizeRoot()...)
		r.Get("/", rh.handleGetKeys)
		r.Post("/", rh.handleCreateDataSet)
	})
	r.Route("/api/{key}/name", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeDataSet()...)
		r.Get("/", rh.handleGetDataSetName)
		r.Put("/", rh.handlePutDataSetName)
	})
	r.Route("/api/{key}", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeDataSet()...)
		r.Get("/", rh.handleGetData)
		r.Put("/", rh.handlePutData)
	})
	r.Route("/api/{key}/{startTime}", func(r chi.Router) {
		r.Use(makeContextSaver(StartTimeKey, "startTime"))
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeDataSet()...)
		r.Get("/", rh.handleGetDataWithStart)
		r.Put("/", rh.handlePutDataWithStart)
		r.Delete("/", rh.handleDeleteData)
	})
	r.Route("/api/{key}/schema", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeDataSet()...)
		r.Get("/", rh.handleGetSchema)
		r.Put("/", rh.handlePutSchema)
	})
	r.Route("/api/{key}/schema/{startTime}", func(r chi.Router) {
		r.Use(makeContextSaver(StartTimeKey, "startTime"))
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeDataSet()...)
		r.Get("/", rh.handleGetSchemaWithStart)
		r.Put("/", rh.handlePutSchemaWithStart)
		r.Delete("/", rh.handleDeleteSchema)
	})
	r.Route("/api/admin/{key}", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeAdmin()...)
		r.Delete("/", rh.handleDeleteDataSet)
	})
	r.Route("/api/admin/{key}/start-times", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeAdmin()...)
		r.Get("/", rh.handleGetStartTimes)
	})
	r.Route("/api/admin/{key}/creation-times", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeAdmin()...)
		r.Get("/", rh.handleGetCreationTimes)
	})
	r.Route("/api/admin/{key}/restrictions", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeAdmin()...)
		r.Get("/", rh.handleGetRestriction)
		r.Put("/", rh.handlePutRestriction)
	})
	r.Route("/api/admin/{key}/subscribe", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeAdmin()...)
		r.Post("/", rh.handlePutSubscription)
	})
	r.Route("/api/admin/{key}/unsubscribe", func(r chi.Router) {
		r.Use(makeContextSaver(KeyKey, "key"))
		r.Use(m.AuthorizeAdmin()...)
		r.Post("/", rh.handleDeleteSubscription)
	})
	return rh
}

func (rh *RestHandler) handleGetKeys(w http.ResponseWriter, r *http.Request) {
	result, err := rh.storer.GetKeys()
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleCreateDataSet(w http.ResponseWriter, r *http.Request) {
	var payload DataSetPayload
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.CreateDataSet(payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleGetDataSetName(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	queryTime := r.URL.Query().Get("time")
	result, err := rh.storer.GetDataSetName(key, time)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handlePutDataSetName(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	queryTime := r.URL.Query().Get("time")
	var payload NamePayload
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.PutDataSetName(key, time, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleGetData(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	queryTime := r.URL.Query().Get("time")
	result, err := rh.storer.GetData(key, time)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handlePutData(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	queryTime := r.URL.Query().Get("time")
	var payload DataPayload
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.PutData(key, time, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleGetDataWithStart(w http.ResponseWriter, r *http.Request) {
	startTime := r.Context().Value(StartTimeKey).(string)
	key := r.Context().Value(KeyKey).(string)
	result, err := rh.storer.GetDataWithStart(startTime, key)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handlePutDataWithStart(w http.ResponseWriter, r *http.Request) {
	startTime := r.Context().Value(StartTimeKey).(string)
	key := r.Context().Value(KeyKey).(string)
	var payload DataPayload
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.PutDataWithStart(startTime, key, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleDeleteData(w http.ResponseWriter, r *http.Request) {
	startTime := r.Context().Value(StartTimeKey).(string)
	key := r.Context().Value(KeyKey).(string)
	if err := rh.storer.DeleteData(startTime, key); err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.NoContent(w, r)
}

func (rh *RestHandler) handleGetSchema(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	queryTime := r.URL.Query().Get("time")
	result, err := rh.storer.GetSchema(key, time)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handlePutSchema(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	queryTime := r.URL.Query().Get("time")
	var payload SchemaPayload
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.PutSchema(key, time, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleGetSchemaWithStart(w http.ResponseWriter, r *http.Request) {
	startTime := r.Context().Value(StartTimeKey).(string)
	key := r.Context().Value(KeyKey).(string)
	result, err := rh.storer.GetSchemaWithStart(startTime, key)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handlePutSchemaWithStart(w http.ResponseWriter, r *http.Request) {
	startTime := r.Context().Value(StartTimeKey).(string)
	key := r.Context().Value(KeyKey).(string)
	var payload SchemaPayload
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.PutSchemaWithStart(startTime, key, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleDeleteSchema(w http.ResponseWriter, r *http.Request) {
	startTime := r.Context().Value(StartTimeKey).(string)
	key := r.Context().Value(KeyKey).(string)
	if err := rh.storer.DeleteSchema(startTime, key); err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.NoContent(w, r)
}

func (rh *RestHandler) handleDeleteDataSet(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	if err := rh.storer.DeleteDataSet(key); err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.NoContent(w, r)
}

func (rh *RestHandler) handleGetStartTimes(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	result, err := rh.storer.GetStartTimes(key)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleGetCreationTimes(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	result, err := rh.storer.GetCreationTimes(key)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleGetRestriction(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	result, err := rh.storer.GetRestriction(key)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handlePutRestriction(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	var payload Restriction
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.PutRestriction(key, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.JSON(w, r, result)
}

func (rh *RestHandler) handlePutSubscription(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	var payload Subscription
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.PutSubscription(key, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, result)
}

func (rh *RestHandler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request) {
	key := r.Context().Value(KeyKey).(string)
	var payload Subscription
	if err := decodeJSON(r.Body, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := rh.storer.DeleteSubscription(key, payload)
	if err != nil {
		http.Error(w, err.Error(), getStatus(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, result)
}
`

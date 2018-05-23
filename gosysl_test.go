package gosysl

import (
	"io/ioutil"
	"testing"

	"github.com/anz-bank/gosysl/pb"
	"github.com/golang/protobuf/proto"
	testifyAssert "github.com/stretchr/testify/assert"
)

var expectedRest = `package mypkg

// Data holds JSON data valid from StartTime created at CreationTime
type Data struct {
	StartTime    string ` + "`" + `json:"start-time"` + "`" + `
	JSONData     string ` + "`" + `json:"data"` + "`" + `
	CreationTime string ` + "`" + `json:"creation-time"` + "`" + `
}

// Schema holds JSON Schema to validate Data against, a name for key creation.
type Schema struct {
	StartTime    string ` + "`" + `json:"start-time"` + "`" + `
	JSONSchema   string ` + "`" + `json:"schema"` + "`" + `
	CreationTime string ` + "`" + `json:"creation-time"` + "`" + `
}

// Subscription holds external endpoint values for change notification
type Subscription struct {
	URL          string ` + "`" + `json:"url"` + "`" + `
	SecreteToken string ` + "`" + `json:"secrete-token"` + "`" + `
}

// Restriction contains scope access restriction and frozen times for schema and data.
type Restriction struct {
	SchemaFrozenUntil string   ` + "`" + `json:"schema-frozen-until"` + "`" + `
	DataFrozenUntil   string   ` + "`" + `json:"data-frozen-until"` + "`" + `
	ReadScopes        []string ` + "`" + `json:"read-scopes"` + "`" + `
	ReadWriteScopes   []string ` + "`" + `json:"read-write-scopes"` + "`" + `
	AdminScopes       []string ` + "`" + `json:"admin-scopes"` + "`" + `
}

// Keys is JSON result type for getKeys in REST API
type Keys struct {
	Keys []string ` + "`" + `json:"keys"` + "`" + `
}

// Key is JSON result type for createDataSet in REST API
type Key struct {
	Key string ` + "`" + `json:"key"` + "`" + `
}

// KeyName is JSON result type for get and put dataDetNamre in REST API
type KeyName struct {
	Key  string ` + "`" + `json:"key"` + "`" + `
	Name string ` + "`" + `json:"name"` + "`" + `
}

// Times contains schema and data times, used to get StartTimes for both
type Times struct {
	Data   []string ` + "`" + `json:"data-times"` + "`" + `
	Schema []string ` + "`" + `json:"schema"` + "`" + `
}

// CreationStartTime contains start and creation time for a schema or data snapshot
type CreationStartTime struct {
	CreationTime string ` + "`" + `json:"creation-time"` + "`" + `
	StartTime    string ` + "`" + `json:"start-time"` + "`" + `
}

// CreationTimes contains schema and data times maps, used to StartTime to CreationTims
type CreationTimes struct {
	Data   map[string]CreationStartTime ` + "`" + `json:"data-time-map"` + "`" + `
	Schema map[string]CreationStartTime ` + "`" + `json:"schema"` + "`" + `
}

// DataSetPayload is JSON payload on REST API request to create new data set
type DataSetPayload struct {
	Name         string      ` + "`" + `json:"name"` + "`" + `
	StartTimeStr string      ` + "`" + `json:"start-time"` + "`" + `
	JSONSchema   interface{} ` + "`" + `json:"schema"` + "`" + `
}

// NamePayload is JSON payload on REST API request to update data set name
type NamePayload struct {
	Name string ` + "`" + `json:"name"` + "`" + `
}

// DataPayload is JSON payload on REST API request to update data
type DataPayload struct {
	Data interface{} ` + "`" + `json:"data"` + "`" + `
}

// SchemaPayload is JSON payload on REST API request to update schema
type SchemaPayload struct {
	Schema interface{} ` + "`" + `json:"schema"` + "`" + `
}

// UpdateEvent holds all information necessary to post to subscribes
type UpdateEvent struct {
	Key       string      ` + "`" + `json:"key"` + "`" + `
	StartTime string      ` + "`" + `json:"start-time"` + "`" + `
	Data      interface{} ` + "`" + `json:"data"` + "`" + `
	Schema    interface{} ` + "`" + `json:"schema"` + "`" + `
	Deleted   bool        ` + "`" + `json:"deleted"` + "`" + `
}
`

func TestEnd2End(tt *testing.T) {
	assert := testifyAssert.New(tt)
	data, err := ioutil.ReadFile("example/example.pb")
	assert.NoError(err)
	module := &pb.Module{}
	err = proto.Unmarshal(data, module)
	assert.NoError(err)
	result, err := Generate(module, "mypkg")
	assert.NoError(err)
	assert.Equal(expectedRest, result.Interface)

	// failing gofmt
	_, err = Generate(module, "BAD PACKAGE NAME")
	assert.Error(err)
}

func TestGetPackage(tt *testing.T) {
	assert := testifyAssert.New(tt)

	assert.Equal("x", GetPackage("x"))
	assert.Equal("y", GetPackage("x/y"))
}

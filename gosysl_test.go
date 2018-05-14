package gosysl

import (
	"io/ioutil"
	"testing"

	"github.com/anz-bank/gosysl/pb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

var expectedRest = `package mypkg

type Data struct {
	StartTime    string
	JSONData     string
	CreationTime string
}
type Schema struct {
	StartTime    string
	JSONSchema   string
	CreationTime string
}
type Subscription struct {
	URL          string
	SecreteToken string
}
type Restriction struct {
	SchemaFrozenUntil string
	DataFrozenUntil   string
	ReadScopes        []string
	ReadWriteScopes   []string
	AdminScopes       []string
}
type Keys struct {
	Keys []string
}
type Key struct {
	Key string
}
type KeyName struct {
	Key  string
	Name string
}
type Times struct {
	Data   []string
	Schema []string
}
type CreationStartTime struct {
	CreationTime string
	StartTime    string
}
type CreationTimes struct {
	Data   map[string]CreationStartTime
	Schema map[string]CreationStartTime
}
type DataSetPayload struct {
	Name         string
	StartTimeStr string
	JSONSchema   interface{}
}
type NamePayload struct {
	Name string
}
type DataPayload struct {
	Data interface{}
}
type SchemaPayload struct {
	Schema interface{}
}
type UpdateEvent struct {
	Key       string
	StartTime string
	Data      interface{}
	Schema    interface{}
	Deleted   bool
}
`

func TestEnd2End(tt *testing.T) {
	assert := assert.New(tt)

	data, err := ioutil.ReadFile("example/example.pb")
	assert.NoError(err)
	module := &pb.Module{}
	err = proto.Unmarshal(data, module)
	assert.NoError(err)
	result, err := Generate(module, "mypkg")
	assert.NoError(err)
	assert.Equal(expectedRest, result.Interface)
}

func TestGetPackage(tt *testing.T) {
	assert := assert.New(tt)

	assert.Equal("x", GetPackage("x"))
	assert.Equal("y", GetPackage("x/y"))
}

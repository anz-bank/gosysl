package gosysl

import (
	"bytes"
	"testing"

	"github.com/anz-bank/gosysl/pb"
	testifyAssert "github.com/stretchr/testify/assert"
)

func TestWriteMethodName(tt *testing.T) {
	assert := testifyAssert.New(tt)
	var tests = []struct {
		input    string
		expected string
	}{
		{"GET /api", "GetApi"},
		{"GET /api/admin/{key}/creation-times", "GetApiAdminKeyCreationTimes"},
	}

	for _, t := range tests {
		ep := &pb.Endpoint{Name: t.input}
		assert.Equal(t.expected, GetMethodName(ep))
	}
}

func TestBadInput(tt *testing.T) {
	assert := testifyAssert.New(tt)

	t := &pb.Type{
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 1}},
	}
	qp := &pb.Endpoint_RestParams_QueryParam{Type: t}
	qpSlice := []*pb.Endpoint_RestParams_QueryParam{qp}
	rp := &pb.Endpoint_RestParams{QueryParam: qpSlice}
	ep := &pb.Endpoint{RestParams: rp}
	_, err := getParams(ep)
	assert.Error(err)

	_, err = getReturnTypes(ep)
	assert.Error(err)

	w := &bytes.Buffer{}
	assert.Error(writeMethod(w, ep))
	ep.RestParams.QueryParam = nil
	assert.Error(writeMethod(w, ep))

	app := &pb.Application{
		Endpoints: map[string]*pb.Endpoint{"x": ep},
	}
	assert.Error(WriteInterface(w, app, []string{"x"}))

	_, err = genInterfaceFile(app, []string{"x"}, "pkg")
	assert.Error(err)

	app.Endpoints = nil
	a := &pb.Attribute{
		Attribute: &pb.Attribute_S{S: "Bad Go interface identifier"},
	}
	attrs := map[string]*pb.Attribute{"interface": a}
	app.Attrs = attrs
	assert.NoError(WriteInterface(w, app, nil))
	module := &pb.Module{
		Apps: map[string]*pb.Application{"x": app},
	}
	_, err = Generate(module, "x")
	assert.Error(err)
}

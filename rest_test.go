package gosysl

import (
	"testing"

	"github.com/anz-bank/gosysl/pb"
	testifyAssert "github.com/stretchr/testify/assert"
)

func TestRestCornerCases(tt *testing.T) {
	assert := testifyAssert.New(tt)

	attr := &pb.Attribute{
		Attribute: &pb.Attribute_S{
			S: "Bad Go Func Name",
		},
	}
	ret := &pb.Statement{
		Stmt: &pb.Statement_Action{Action: &pb.Action{Action: "return"}},
	}
	ep := &pb.Endpoint{
		Attrs: map[string]*pb.Attribute{"middleware": attr},
		Stmt:  []*pb.Statement{ret},
		Name:  "endpoint",
	}
	app := &pb.Application{
		Endpoints: map[string]*pb.Endpoint{"ep": ep},
	}
	_, err := GenMiddleware(app, []string{"ep"})
	assert.Error(err)

	module := &pb.Module{
		Apps: map[string]*pb.Application{"app": app},
	}
	_, err = Generate(module, "pkg")
	assert.Error(err)

	_, err = genMiddlewareFile(&pb.Application{}, nil, "BAD PACKAGE")
	assert.Error(err)

}

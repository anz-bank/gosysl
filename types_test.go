package gosysl

import (
	"bytes"
	"go/format"
	"testing"

	"github.com/anz-bank/gosysl/pb"
	testifyAssert "github.com/stretchr/testify/assert"
)

var primitiveTests = []struct {
	input    pb.Type_Primitive
	expected string
}{
	{pb.Type_BOOL, "bool"},
	{pb.Type_ANY, "interface{}"},
	{pb.Type_INT, "int"},
	{pb.Type_STRING, "string"},
	{pb.Type_FLOAT, "float64"},
	{pb.Type_DECIMAL, "float64"},
	{pb.Type_EMPTY, "nil"},
	{pb.Type_BYTES, "[]byte"},
	{pb.Type_DATE, "time.Time"},
	{pb.Type_DATETIME, "time.Time"},
	{pb.Type_STRING, "string"},
}
var primitiveErrTests = []pb.Type_Primitive{
	pb.Type_XML,
	pb.Type_UUID,
}

func TestGetPrimitiveType(tt *testing.T) {
	assert := testifyAssert.New(tt)
	for _, t := range primitiveTests {
		str, err := GetPrimitiveType(&t.input)
		assert.NoError(err)
		assert.Equal(t.expected, str)
	}
	for _, tp := range primitiveErrTests {
		_, err := GetPrimitiveType(&tp)
		assert.Error(err)
	}
}

func TestGetSimpleType(tt *testing.T) {
	assert := testifyAssert.New(tt)
	for _, t := range primitiveTests {
		primitive := &pb.Type_Primitive_{Primitive: t.input}
		pt := &pb.Type{Type: primitive}
		str, err := GetSimpleType(pt)
		assert.NoError(err)
		assert.Equal(t.expected, str)
	}
	for _, tp := range primitiveErrTests {
		primitive := &pb.Type_Primitive_{Primitive: tp}
		pt := &pb.Type{Type: primitive}
		_, err := GetSimpleType(pt)
		assert.Error(err)
	}

	// test type refs
	var refTests = []struct {
		input    string
		expected string
	}{
		{"MyType", "MyType"},
		{"x", "x"},
		{"map of string:int", "map[string]int"},
	}
	for _, t := range refTests {
		ref := &pb.Scope{Path: []string{t.input}}
		typeRef := &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}}
		pt := &pb.Type{Type: typeRef}
		str, err := GetSimpleType(pt)
		assert.NoError(err)
		assert.Equal(t.expected, str)
	}
	errRefTests := [][]string{
		{"MyType", "MyType"},
		{"map of string:int:bool"},
	}
	for _, t := range errRefTests {
		ref := &pb.Scope{Path: t}
		typeRef := &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}}
		pt := &pb.Type{Type: typeRef}
		_, err := GetSimpleType(pt)
		assert.Error(err)
	}

	strType := &pb.Type_Primitive_{Primitive: pb.Type_STRING}
	list := &pb.Type_List_{List: &pb.Type_List{Type: &pb.Type{Type: strType}}}
	pt := &pb.Type{Type: list}
	_, err := GetSimpleType(pt)
	assert.Error(err)
}

func TestSplitUppercase(tt *testing.T) {
	assert := testifyAssert.New(tt)

	assert.Equal([]string{"a"}, SplitUppercase("a"))
	assert.Equal([]string{"abc"}, SplitUppercase("abc"))
	assert.Equal([]string{"abc", "xyz"}, SplitUppercase("AbcXyz"))
	assert.Equal([]string{"a", "b", "c"}, SplitUppercase("ABC"))
	assert.Equal([]string{"a", "b", "c"}, SplitUppercase("aBC"))
}

func TestWriteStruct(tt *testing.T) {
	assert := testifyAssert.New(tt)

	attrDefs := map[string]*pb.Type{}
	typeTuple := &pb.Type_Tuple_{Tuple: &pb.Type_Tuple{AttrDefs: attrDefs}}
	ttype := &pb.Type{Type: typeTuple}

	attrDefs["Data"] = &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_ANY},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 1}},
	}

	w := &bytes.Buffer{}
	assert.NoError(WriteStruct(w, "DataPayload", ttype, ""))
	expectedSrc := `type DataPayload struct {
			Data interface{} ` + "`json:\"Data\"`" + `
		}` + "\n"
	expected, _ := format.Source([]byte(expectedSrc))
	actual, err := format.Source(w.Bytes())
	assert.NoError(err)
	assert.Equal(string(expected), string(actual))

	attrDefs["LastName"] = &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_STRING},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 2}},
	}

	w = &bytes.Buffer{}
	assert.NoError(WriteStruct(w, "DataPayload", ttype, "-"))
	expectedSrc = `type DataPayload struct {
			Data interface{} ` + "`json:\"data\"`" + `
			LastName string  ` + "`json:\"last-name\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	actual, err = format.Source(w.Bytes())
	assert.NoError(err)
	assert.Equal(string(expected), string(actual))

	ref := &pb.Scope{Path: []string{"MyType"}}
	attrDefs["MyField"] = &pb.Type{
		Type:          &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 3}},
	}
	w = &bytes.Buffer{}
	assert.NoError(WriteStruct(w, "DataPayload", ttype, "__"))
	expectedSrc = `type DataPayload struct {
			Data interface{} ` + "`json:\"data\"`" + `
			LastName string  ` + "`json:\"last__name\"`" + `
			MyField MyType   ` + "`json:\"my__field\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	actual, err = format.Source(w.Bytes())
	assert.NoError(err)
	assert.Equal(string(expected), string(actual))

	strType := &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_STRING},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 100}},
	}
	typeList := &pb.Type_List_{List: &pb.Type_List{Type: strType}}
	attrDefs["StringList"] = &pb.Type{
		Type: typeList,
	}
	w = &bytes.Buffer{}
	assert.NoError(WriteStruct(w, "DataPayload", ttype, ""))
	expectedSrc = `type DataPayload struct {
			Data interface{}   ` + "`json:\"Data\"`" + `
			LastName string    ` + "`json:\"LastName\"`" + `
			MyField MyType     ` + "`json:\"MyField\"`" + `
			StringList []string` + "`json:\"StringList\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	actual, err = format.Source(w.Bytes())
	assert.NoError(err)
	assert.Equal(string(expected), string(actual))

	refType := &pb.Type{
		Type:          &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 10}},
	}
	typeList = &pb.Type_List_{List: &pb.Type_List{Type: refType}}
	attrDefs["MyList"] = &pb.Type{
		Type: typeList,
	}
	w = &bytes.Buffer{}
	assert.NoError(WriteStruct(w, "DataPayload", ttype, "-"))
	expectedSrc = `type DataPayload struct {
			Data interface{}   ` + "`json:\"data\"`" + `
			LastName string    ` + "`json:\"last-name\"`" + `
			MyField MyType     ` + "`json:\"my-field\"`" + `
			MyList []MyType    ` + "`json:\"my-list\"`" + `
			StringList []string` + "`json:\"string-list\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	actual, err = format.Source(w.Bytes())
	assert.NoError(err)
	assert.Equal(string(expected), string(actual))

	strType2 := &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_STRING},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 20}},
	}
	attrDefs["StringSet"] = &pb.Type{
		Type: &pb.Type_Set{Set: strType2},
	}
	w = &bytes.Buffer{}
	assert.NoError(WriteStruct(w, "YYY", ttype, ""))
	expectedSrc = `type YYY struct {
			Data interface{}                 ` + "`json:\"Data\"`" + `
			LastName string                  ` + "`json:\"LastName\"`" + `
			MyField MyType                   ` + "`json:\"MyField\"`" + `
			MyList []MyType                  ` + "`json:\"MyList\"`" + `
			StringSet map[string]interface{} ` + "`json:\"StringSet\"`" + `
			StringList []string              ` + "`json:\"StringList\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	actual, err = format.Source(w.Bytes())
	assert.NoError(err)
	assert.Equal(string(expected), string(actual))
}

func TestGentypesCornerCases(tt *testing.T) {
	assert := testifyAssert.New(tt)

	module := &pb.Module{}
	_, err := Generate(module, "")
	assert.Error(err)

	_, err = GetTypeLine(&pb.Type{})
	assert.Error(err)

	types := map[string]*pb.Type{"x": {}}
	_, err = NamesSortedBySourceContext(types)
	assert.Error(err)
	w := &bytes.Buffer{}
	assert.Error(WriteTypes(w, &pb.Application{Types: types}))
	module.Apps = map[string]*pb.Application{"x": {Types: types}}
	_, err = Generate(module, "")
	assert.Error(err)

	assert.Error(WriteStruct(w, "x", &pb.Type{}, "-"))

	attrDefs := map[string]*pb.Type{}
	typeTuple := &pb.Type_Tuple_{Tuple: &pb.Type_Tuple{AttrDefs: attrDefs}}
	ttype := &pb.Type{Type: typeTuple}
	attrDefs["Data"] = &pb.Type{}
	assert.Error(WriteStruct(w, "DataPayload", ttype, ""))

	ref := &pb.Scope{Path: []string{"1", "2"}}
	attrDefs["Data"] = &pb.Type{
		Type:          &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 3}},
	}
	assert.Error(WriteStruct(w, "DataPayload", ttype, ""))
	types = map[string]*pb.Type{"x": ttype}
	assert.Error(WriteTypes(w, &pb.Application{Types: types}))

	_, _, err = GetType(&pb.Type{})
	assert.Error(err)

	assert.Error(WriteStructField(w, "", &pb.Type{}, ""))

}

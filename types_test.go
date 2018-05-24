package gosysl

import (
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

func TestGenPrimitiveType(tt *testing.T) {
	assert := testifyAssert.New(tt)
	for _, t := range primitiveTests {
		str, err := GenPrimitiveType(&t.input)
		assert.NoError(err)
		assert.Equal(t.expected, str)
	}
	for _, tp := range primitiveErrTests {
		_, err := GenPrimitiveType(&tp)
		assert.Error(err)
	}
}

func TestGenSimpleType(tt *testing.T) {
	assert := testifyAssert.New(tt)
	for _, t := range primitiveTests {
		primitive := &pb.Type_Primitive_{Primitive: t.input}
		pt := &pb.Type{Type: primitive}
		str, err := GenSimpleType(pt)
		assert.NoError(err)
		assert.Equal(t.expected, str)
	}
	for _, tp := range primitiveErrTests {
		primitive := &pb.Type_Primitive_{Primitive: tp}
		pt := &pb.Type{Type: primitive}
		_, err := GenSimpleType(pt)
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
		str, err := GenSimpleType(pt)
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
		_, err := GenSimpleType(pt)
		assert.Error(err)
	}

	strType := &pb.Type_Primitive_{Primitive: pb.Type_STRING}
	list := &pb.Type_List_{List: &pb.Type_List{Type: &pb.Type{Type: strType}}}
	pt := &pb.Type{Type: list}
	_, err := GenSimpleType(pt)
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

func TestGenStruct(tt *testing.T) {
	assert := testifyAssert.New(tt)

	attrDefs := map[string]*pb.Type{}
	typeTuple := &pb.Type_Tuple_{Tuple: &pb.Type_Tuple{AttrDefs: attrDefs}}
	ttype := &pb.Type{Type: typeTuple}

	attrDefs["Data"] = &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_ANY},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 1}},
	}

	str, err := GenStruct("DataPayload", ttype, "")
	assert.NoError(err)
	expectedSrc := `type DataPayload struct {
			Data interface{} ` + "`json:\"Data\"`" + `
		}` + "\n"
	expected, _ := format.Source([]byte(expectedSrc))
	assert.Equal(string(expected), str)

	attrDefs["LastName"] = &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_STRING},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 2}},
	}
	str, err = GenStruct("DataPayload", ttype, "-")
	assert.NoError(err)
	expectedSrc = `type DataPayload struct {
			Data interface{} ` + "`json:\"data\"`" + `
			LastName string  ` + "`json:\"last-name\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	assert.Equal(string(expected), str)

	ref := &pb.Scope{Path: []string{"MyType"}}
	attrDefs["MyField"] = &pb.Type{
		Type:          &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 3}},
	}
	str, err = GenStruct("DataPayload", ttype, "__")
	assert.NoError(err)
	expectedSrc = `type DataPayload struct {
			Data interface{} ` + "`json:\"data\"`" + `
			LastName string  ` + "`json:\"last__name\"`" + `
			MyField MyType   ` + "`json:\"my__field\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	assert.Equal(string(expected), str)

	strType := &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_STRING},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 100}},
	}
	typeList := &pb.Type_List_{List: &pb.Type_List{Type: strType}}
	attrDefs["StringList"] = &pb.Type{
		Type: typeList,
	}
	str, err = GenStruct("DataPayload", ttype, "")
	assert.NoError(err)
	expectedSrc = `type DataPayload struct {
			Data interface{}   ` + "`json:\"Data\"`" + `
			LastName string    ` + "`json:\"LastName\"`" + `
			MyField MyType     ` + "`json:\"MyField\"`" + `
			StringList []string` + "`json:\"StringList\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	assert.Equal(string(expected), str)

	refType := &pb.Type{
		Type:          &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 10}},
	}
	typeList = &pb.Type_List_{List: &pb.Type_List{Type: refType}}
	attrDefs["MyList"] = &pb.Type{
		Type: typeList,
	}
	str, err = GenStruct("DataPayload", ttype, "-")
	assert.NoError(err)
	expectedSrc = `type DataPayload struct {
			Data interface{}   ` + "`json:\"data\"`" + `
			LastName string    ` + "`json:\"last-name\"`" + `
			MyField MyType     ` + "`json:\"my-field\"`" + `
			MyList []MyType    ` + "`json:\"my-list\"`" + `
			StringList []string` + "`json:\"string-list\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	assert.Equal(string(expected), str)

	strType2 := &pb.Type{
		Type:          &pb.Type_Primitive_{Primitive: pb.Type_STRING},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 20}},
	}
	attrDefs["StringSet"] = &pb.Type{
		Type: &pb.Type_Set{Set: strType2},
	}
	str, err = GenStruct("YYY", ttype, "")
	assert.NoError(err)
	expectedSrc = `type YYY struct {
			Data interface{}                 ` + "`json:\"Data\"`" + `
			LastName string                  ` + "`json:\"LastName\"`" + `
			MyField MyType                   ` + "`json:\"MyField\"`" + `
			MyList []MyType                  ` + "`json:\"MyList\"`" + `
			StringSet map[string]interface{} ` + "`json:\"StringSet\"`" + `
			StringList []string              ` + "`json:\"StringList\"`" + `
		}` + "\n"
	expected, _ = format.Source([]byte(expectedSrc))
	assert.Equal(string(expected), str)
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
	_, err = GenTypes(&pb.Application{Types: types})
	assert.Error(err)
	module.Apps = map[string]*pb.Application{"x": {Types: types}}
	_, err = Generate(module, "")
	assert.Error(err)

	_, err = GenStruct("x", &pb.Type{}, "-")
	assert.Error(err)

	attrDefs := map[string]*pb.Type{}
	typeTuple := &pb.Type_Tuple_{Tuple: &pb.Type_Tuple{AttrDefs: attrDefs}}
	ttype := &pb.Type{Type: typeTuple}
	attrDefs["Data"] = &pb.Type{}
	_, err = GenStruct("DataPayload", ttype, "")
	assert.Error(err)

	ref := &pb.Scope{Path: []string{"1", "2"}}
	attrDefs["Data"] = &pb.Type{
		Type:          &pb.Type_TypeRef{TypeRef: &pb.ScopedRef{Ref: ref}},
		SourceContext: &pb.SourceContext{Start: &pb.SourceContext_Location{Line: 3}},
	}
	_, err = GenStruct("DataPayload", ttype, "")
	assert.Error(err)
	types = map[string]*pb.Type{"x": ttype}
	_, err = GenTypes(&pb.Application{Types: types})
	assert.Error(err)

	_, _, err = GenType(&pb.Type{})
	assert.Error(err)

	_, err = GenStructField("", &pb.Type{}, "")
	assert.Error(err)

}

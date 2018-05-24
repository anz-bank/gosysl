package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"unicode"

	"github.com/anz-bank/gosysl/pb"
)

// GenPrimitiveType return Golang type string for Sysl primitive data type
func GenPrimitiveType(tp *pb.Type_Primitive) (string, error) {
	switch *tp {
	case pb.Type_BOOL:
		return "bool", nil
	case pb.Type_ANY:
		return "interface{}", nil
	case pb.Type_INT:
		return "int", nil
	case pb.Type_STRING:
		return "string", nil
	case pb.Type_EMPTY:
		return "nil", nil
	case pb.Type_FLOAT, pb.Type_DECIMAL:
		return "float64", nil
	case pb.Type_BYTES:
		return "[]byte", nil
	case pb.Type_DATE, pb.Type_DATETIME:
		return "time.Time", nil
	}
	// Type_XML, Type_UUID
	return "", fmt.Errorf("unsupported type primitive %s", tp.String())
}

// GenSimpleType returns Golang type string for non-composite types (no lists and sets)
func GenSimpleType(t *pb.Type) (string, error) {
	if pType := t.GetPrimitive(); pType != pb.Type_NO_Primitive {
		return GenPrimitiveType(&pType)
	}
	if t.GetTypeRef() != nil {
		path := t.GetTypeRef().GetRef().GetPath()
		if len(path) != 1 {
			return "", fmt.Errorf("cannot handle type reference with more than one path")
		}
		str := path[0]
		if !strings.HasPrefix(str, "map of") {
			return str, nil
		}
		str = strings.TrimPrefix(str, "map of")
		str = strings.Trim(str, " ")
		m := strings.Split(str, ":")
		if len(m) != 2 {
			shouldStr := `should be map of KeyType:ValueType`
			return "", fmt.Errorf("bad map definition '%s' (%s)", str, shouldStr)
		}
		return fmt.Sprintf("map[%s]%s", m[0], m[1]), nil
	}
	return "", fmt.Errorf("type %v is neither primitive nor reference", t)

}

// GetLine returns the line for a given Sysl type from its SourceContext
func GetLine(t *pb.Type) (int32, error) {
	if t.GetPrimitive() != pb.Type_NO_Primitive || t.GetTypeRef() != nil {
		return t.GetSourceContext().Start.Line, nil
	}
	if t.GetList() != nil {
		return t.GetList().GetType().GetSourceContext().Start.Line, nil
	}
	if t.GetSet() != nil {
		return t.GetSet().GetSourceContext().Start.Line, nil
	}
	if t.GetTuple() != nil {
		for _, t2 := range t.GetTuple().GetAttrDefs() {
			return GetLine(t2)
		}
	}
	return 0, fmt.Errorf("unknown type %v for getting line", t)
}

//NamesSortedBySourceContext sorts the keys of the input types according to occurrence
// in Sysl definition file, derived from SourceContext
func NamesSortedBySourceContext(types map[string]*pb.Type) ([]string, error) {
	type lineName struct {
		name string
		line int32
	}
	size := len(types)
	lineNames := make([]lineName, size)
	i := 0
	for name, t := range types {
		line, err := GetLine(t)
		if err != nil {
			return nil, err
		}
		lineNames[i] = lineName{name, line}
		i++
	}
	sort.Slice(lineNames, func(i, j int) bool {
		return lineNames[i].line < lineNames[j].line
	})
	result := make([]string, size)
	for i := 0; i < size; i++ {
		result[i] = lineNames[i].name
	}
	return result, nil
}

// SplitUppercase splits into fields at all upper case letters and converts
// fields to lower case
func SplitUppercase(str string) []string {
	result := make([]string, 0, 8)
	idx := make([]int, 0, 8)
	for pos, c := range str {
		if unicode.IsUpper(c) || pos == 0 {
			idx = append(idx, pos)
		}
	}
	idx = append(idx, len(str))
	for i := 0; i < len(idx)-1; i++ {
		result = append(result, strings.ToLower(str[idx[i]:idx[i+1]]))
	}
	return result
}

// GetJSONProperty extracts JSON property name from attribute and defaults to type name
func GetJSONProperty(name string, t *pb.Type, sep string) string {
	if attrOverride, ok := t.Attrs["json"]; ok {
		return attrOverride.GetS()
	}
	if sep == "" {
		return name
	}
	return strings.Join(SplitUppercase(name), sep)
}

// GenStructField creates a single line inside a struct definition
func GenStructField(fieldName string, fieldType *pb.Type, sep string) (string, error) {
	var typeStr string
	var err error
	if fieldType.GetPrimitive() != pb.Type_NO_Primitive || fieldType.GetTypeRef() != nil {
		// Primitive type, reference or map
		if typeStr, err = GenSimpleType(fieldType); err == nil {
			typeStr = fmt.Sprintf("%s %s ", fieldName, typeStr)
		}
	} else if fieldType.GetList() != nil {
		// List
		fieldType = fieldType.GetList().GetType()
		if typeStr, err = GenSimpleType(fieldType); err == nil {
			typeStr = fmt.Sprintf("%s []%s ", fieldName, typeStr)
		}
	} else if fieldType.GetSet() != nil {
		// Set
		fieldType = fieldType.GetSet()
		if typeStr, err = GenSimpleType(fieldType); err == nil {
			typeStr = fmt.Sprintf("%s map[%s]interface{} ", fieldName, typeStr)
		}
	}
	if err != nil {
		return "", err
	}
	if typeStr == "" {
		return "", fmt.Errorf("unknown type %s", fieldType.String())
	}
	jsonProp := GetJSONProperty(fieldName, fieldType, sep)
	return fmt.Sprintf("%s `json:\"%s\"`\n", typeStr, jsonProp), nil
}

//GenStruct creates a Golang `struct` type definition from a Sysl Tuple type definition
func GenStruct(name string, t *pb.Type, jsonSep string) (string, error) {
	if t.GetTuple() == nil {
		return "", fmt.Errorf("top level type has to be Tuple")
	}
	var buffer bytes.Buffer
	if attr, ok := t.Attrs["doc"]; ok {
		buffer.WriteString(fmt.Sprintf("// %s\n", attr.GetS()))
	}
	buffer.WriteString(fmt.Sprintf("type %s struct{\n", name))
	attrDefs := t.GetTuple().GetAttrDefs()

	names, err := NamesSortedBySourceContext(attrDefs)
	if err != nil {
		return "", err
	}

	for _, fieldName := range names {
		var structField string
		structField, err = GenStructField(fieldName, attrDefs[fieldName], jsonSep)
		if err != nil {
			return "", err
		}
		buffer.WriteString(structField)
	}
	buffer.WriteString("}\n")
	b, err := format.Source(buffer.Bytes())
	return string(b), err
}

// GenTypes creates all types definition in SourceContext order for given Sysl
// type definition
func GenTypes(types map[string]*pb.Type, jsonSep string) (string, error) {
	names, err := NamesSortedBySourceContext(types)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	for _, name := range names {
		t := types[name]
		str, err := GenStruct(name, t, jsonSep)
		if err != nil {
			return "", err
		}
		buffer.WriteString(str)
	}
	return buffer.String(), nil
}

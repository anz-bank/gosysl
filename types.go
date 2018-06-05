package gosysl

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/anz-bank/gosysl/pb"
)

// GetPrimitiveType return Golang type string for Sysl primitive data type
func GetPrimitiveType(tp *pb.Type_Primitive) (string, error) {
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

// GetSimpleType returns Golang type string for non-composite types (no lists and sets)
func GetSimpleType(t *pb.Type) (string, error) {
	if pType := t.GetPrimitive(); pType != pb.Type_NO_Primitive {
		return GetPrimitiveType(&pType)
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

// GetTypeLine returns the line for a given Sysl type from its SourceContext
func GetTypeLine(t *pb.Type) (int32, error) {
	if t.GetPrimitive() != pb.Type_NO_Primitive || t.GetTypeRef() != nil {
		return t.SourceContext.Start.Line, nil
	}
	if t.GetList() != nil {
		return t.GetList().GetType().SourceContext.Start.Line, nil
	}
	if t.GetSet() != nil {
		return t.GetSet().SourceContext.Start.Line, nil
	}
	if t.GetTuple() != nil {
		for _, t2 := range t.GetTuple().GetAttrDefs() {
			return GetTypeLine(t2)
		}
	}
	return 0, fmt.Errorf("unknown type %v for getting line", t)
}

//NamesSortedBySourceContext sorts the keys of the input types according to occurrence
// in Sysl definition file, derived from SourceContext
func NamesSortedBySourceContext(types map[string]*pb.Type) ([]string, error) {
	lineNames := make([]LineName, len(types))
	i := 0
	for name, t := range types {
		line, err := GetTypeLine(t)
		if err != nil {
			return nil, err
		}
		lineNames[i] = LineName{name, line}
		i++
	}
	return SortLineNames(lineNames), nil
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

// GetType creates golang type for given sysl type
func GetType(t *pb.Type) (string, *pb.Type, error) {
	var err error
	var typeStr string
	if t.GetPrimitive() != pb.Type_NO_Primitive || t.GetTypeRef() != nil {
		// Primitive type, reference or map
		if typeStr, err = GetSimpleType(t); err == nil {
			return typeStr, t, nil
		}
	} else if t.GetList() != nil {
		// List
		t = t.GetList().GetType()
		if typeStr, err = GetSimpleType(t); err == nil {
			return fmt.Sprintf("[]%s", typeStr), t, nil
		}
	} else if t.GetSet() != nil {
		// Set
		t = t.GetSet()
		if typeStr, err = GetSimpleType(t); err == nil {
			return fmt.Sprintf("map[%s]interface{}", typeStr), t, nil
		}
	}
	if err != nil {
		return "", nil, err
	}
	return "", nil, fmt.Errorf("unknown type %s", t.String())
}

// WriteStructField creates a single line inside a struct definition
func WriteStructField(w io.Writer, fName string, fType *pb.Type, sep string) error {
	fTypeStr, subType, err := GetType(fType)
	if err != nil {
		return err
	}
	jsonProp := GetJSONProperty(fName, subType, sep)
	fmt.Fprintf(w, "%s %s `json:\"%s\"`\n", fName, fTypeStr, jsonProp)
	return nil
}

// WriteStruct creates a Golang `struct` type definition from a Sysl Tuple type definition
func WriteStruct(w io.Writer, name string, t *pb.Type, jsonSep string) error {
	if t.GetTuple() == nil {
		return fmt.Errorf("top level type has to be Tuple")
	}
	if attr, ok := t.Attrs["doc"]; ok {
		fmt.Fprintf(w, "// %s\n", attr.GetS())
	}
	fmt.Fprintf(w, "type %s struct{\n", name)
	attrDefs := t.GetTuple().GetAttrDefs()

	names, err := NamesSortedBySourceContext(attrDefs)
	if err != nil {
		return err
	}

	for _, fieldName := range names {
		if err = WriteStructField(w, fieldName, attrDefs[fieldName], jsonSep); err != nil {
			return err
		}
	}
	fmt.Fprintln(w, "}")
	return nil
}

// WriteTypes creates all types definition in SourceContext order for given Sysl
// type definition
func WriteTypes(w io.Writer, app *pb.Application) error {
	types := app.GetTypes()
	names, err := NamesSortedBySourceContext(types)
	if err != nil {
		return err
	}
	var jsonSep string
	if attr, ok := app.Attrs["json_property_separator"]; ok {
		jsonSep = attr.GetS()
	}

	for _, name := range names {
		t := types[name]
		if err := WriteStruct(w, name, t, jsonSep); err != nil {
			return err
		}
	}
	return nil
}

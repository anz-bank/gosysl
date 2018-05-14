package gosysl

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"

	"github.com/anz-bank/gosysl/pb"
)

// GenPrimitiveType return golang type string for sysl primitive data type
func GenPrimitiveType(tp pb.Type_Primitive) (string, error) {
	switch tp {
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

// GenSimpleType returns golang type string for non-composite types (no lists and sets)
func GenSimpleType(t *pb.Type) (string, error) {
	pType := t.GetPrimitive()
	if pType != pb.Type_NO_Primitive {
		return GenPrimitiveType(pType)
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
			return "", fmt.Errorf("bad map definition '%s', should be map of KeyType:ValueType", str)
		}
		return fmt.Sprintf("map[%s]%s", m[0], m[1]), nil
	}
	return "", fmt.Errorf("type %v is neither primitive nor reference", t)

}

// GetLine returns the line for a given sysl type form its SourceContext
func GetLine(t *pb.Type) (int32, error) {
	if t.GetPrimitive() != pb.Type_NO_Primitive || t.GetTypeRef() != nil {
		return t.GetSourceContext().Start.Line, nil
	} else if t.GetList() != nil {
		return t.GetList().GetType().GetSourceContext().Start.Line, nil
	} else if t.GetSet() != nil {
		return t.GetSet().GetSourceContext().Start.Line, nil
	} else if t.GetTuple() != nil {
		for _, t2 := range t.GetTuple().GetAttrDefs() {
			return GetLine(t2)
		}
	}
	return 0, fmt.Errorf("unknwon type %v for getting line", t)
}

//NamesSortedBySourceContext sorts the keys of the imput types according to occurance
// in sysl definiton file, derived from SourceContext
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

//GenStruct creates a golang struct type defintion from a sysl Tuple type definiotn
func GenStruct(name string, t *pb.Type) (string, error) {
	if t.GetTuple() == nil {
		return "", fmt.Errorf("top level type has to be Tuple")
	}
	buffer := bytes.NewBufferString(fmt.Sprintf("type %s struct{\n", name))
	attrDefs := t.GetTuple().GetAttrDefs()

	names, err := NamesSortedBySourceContext(attrDefs)
	if err != nil {
		return "", err
	}

	for _, fieldName := range names {
		//TODO add JSON fields, docs
		fieldType := attrDefs[fieldName]
		var err error
		var typeStr string
		if fieldType.GetPrimitive() != pb.Type_NO_Primitive || fieldType.GetTypeRef() != nil {
			// Primitive type, reference or map
			if typeStr, err = GenSimpleType(fieldType); err == nil {
				buffer.WriteString(fmt.Sprintf("%s %s\n", fieldName, typeStr))
			}
		} else if fieldType.GetList() != nil {
			// List
			if typeStr, err = GenSimpleType(fieldType.GetList().GetType()); err == nil {
				buffer.WriteString(fmt.Sprintf("%s []%s\n", fieldName, typeStr))
			}
		} else if fieldType.GetSet() != nil {
			// Set
			if typeStr, err = GenSimpleType(fieldType.GetSet()); err == nil {
				buffer.WriteString(fmt.Sprintf("%s map[%s]interface{}\n", fieldName, typeStr))
			}
		}
		if err != nil {
			return "", err
		}
	}
	buffer.WriteString("}\n")
	b, err := format.Source(buffer.Bytes())
	return string(b), err
}

// GenTypes creates all types defintions in SourceContext order for given Sysl type defintion
func GenTypes(types map[string]*pb.Type) (string, error) {
	names, err := NamesSortedBySourceContext(types)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	for _, name := range names {
		t := types[name]
		str, err := GenStruct(name, t)
		if err != nil {
			return "", err
		}
		buffer.WriteString(str)
	}
	return buffer.String(), nil
}

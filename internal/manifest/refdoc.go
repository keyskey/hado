package manifest

import (
	"reflect"
	"slices"
	"strings"
)

type pathEntry struct {
	path      string
	typeDesc  string
	omitempty bool
}

func manifestYAMLPaths() ([]pathEntry, error) {
	var out []pathEntry
	walkManifestPaths(reflect.TypeOf(Manifest{}), "", false, &out)
	slices.SortFunc(out, func(a, b pathEntry) int {
		return strings.Compare(a.path, b.path)
	})
	return out, nil
}

func walkManifestPaths(t reflect.Type, prefix string, parentOmitempty bool, out *[]pathEntry) {
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		y := yamlKey(sf.Tag)
		if y == "" || y == "-" {
			continue
		}
		omitempty := parentOmitempty || yamlOmitempty(sf.Tag)
		path := y
		if prefix != "" {
			path = prefix + "." + y
		}
		ft := sf.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		switch ft.Kind() {
		case reflect.Struct:
			typeDesc := "object"
			*out = append(*out, pathEntry{path: path, typeDesc: typeDesc, omitempty: omitempty})
			walkManifestPaths(ft, path, omitempty, out)
		case reflect.Slice:
			elem := ft.Elem()
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			if elem.Kind() == reflect.Struct {
				*out = append(*out, pathEntry{path: path, typeDesc: "array of object", omitempty: omitempty})
				walkManifestPaths(elem, path, omitempty, out)
			} else {
				*out = append(*out, pathEntry{
					path:      path,
					typeDesc:  "array of " + scalarTypeName(elem),
					omitempty: omitempty,
				})
			}
		default:
			*out = append(*out, pathEntry{
				path:      path,
				typeDesc:  scalarTypeName(ft),
				omitempty: omitempty,
			})
		}
	}
}

func yamlKey(tag reflect.StructTag) string {
	t := tag.Get("yaml")
	if i := strings.Index(t, ","); i >= 0 {
		return t[:i]
	}
	return t
}

func yamlOmitempty(tag reflect.StructTag) bool {
	return strings.Contains(tag.Get("yaml"), "omitempty")
}

func scalarTypeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "boolean"
	default:
		return t.String()
	}
}

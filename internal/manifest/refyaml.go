package manifest

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// WriteManifestReferenceYAML writes a commented YAML manifest enumerating every supported field.
// Descriptions come from manifestYAMLDoc; shape follows types.go. The output is valid YAML for parseManifestBytes.
func WriteManifestReferenceYAML(w io.Writer) error {
	paths, err := manifestYAMLPaths()
	if err != nil {
		return err
	}
	typeByPath := make(map[string]string, len(paths))
	for _, p := range paths {
		typeByPath[p.path] = p.typeDesc
	}
	var b strings.Builder
	writeYAMLFileHeader(&b)
	if err := emitStructYAML(reflect.TypeOf(Manifest{}), "", 0, false, typeByPath, &b); err != nil {
		return err
	}
	_, err = io.WriteString(w, b.String())
	return err
}

func writeYAMLFileHeader(b *strings.Builder) {
	b.WriteString("# HADO manifest reference — GENERATED FILE; do not edit by hand.\n")
	b.WriteString("# Regenerate: make gen-manifest-doc  (or: go run ./cmd/hado manifest doc --out docs/hado.manifest.reference.yaml)\n")
	b.WriteString("# Types: internal/manifest/types.go  Descriptions: internal/manifest/field_docs.go\n")
	b.WriteString("\n")
}

func emitStructYAML(t reflect.Type, pathPrefix string, level int, parentOmitempty bool, typeByPath map[string]string, b *strings.Builder) error {
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("emitStructYAML: not a struct")
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
		childParentOmit := parentOmitempty || yamlOmitempty(sf.Tag)
		_ = childParentOmit

		path := y
		if pathPrefix != "" {
			path = pathPrefix + "." + y
		}
		doc := manifestYAMLDoc[path]
		if strings.TrimSpace(doc) == "" {
			return fmt.Errorf("manifestYAMLDoc missing description for path %q", path)
		}
		writeYAMLCommentBlock(level, docWithLogicalType(path, doc, typeByPath), b)

		ft := sf.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		pad := indentSpaces(level)

		switch ft.Kind() {
		case reflect.Struct:
			b.WriteString(pad)
			b.WriteString(y)
			b.WriteString(":\n")
			if err := emitStructYAML(ft, path, level+1, childParentOmit, typeByPath, b); err != nil {
				return err
			}
		case reflect.Slice:
			elem := ft.Elem()
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			b.WriteString(pad)
			b.WriteString(y)
			if elem.Kind() == reflect.Struct {
				b.WriteString(":\n")
				if err := emitSliceOfStructYAML(elem, path, level, typeByPath, b); err != nil {
					return err
				}
			} else {
				b.WriteString(": ")
				if err := writeYAMLScalarNoNewline(ft, path, b); err != nil {
					return err
				}
				b.WriteString("\n")
			}
		default:
			b.WriteString(pad)
			b.WriteString(y)
			b.WriteString(": ")
			if err := writeYAMLScalarNoNewline(ft, path, b); err != nil {
				return err
			}
			b.WriteString("\n")
		}
	}
	return nil
}

func emitSliceOfStructYAML(elem reflect.Type, elemPathPrefix string, level int, typeByPath map[string]string, b *strings.Builder) error {
	padList := indentSpaces(level + 1)
	padCont := indentSpaces(level + 2)

	fields, err := exportedYAMLFields(elem)
	if err != nil {
		return err
	}
	if len(fields) == 0 {
		return fmt.Errorf("slice %s: no exported yaml fields", elemPathPrefix)
	}

	for fi, sf := range fields {
		y := yamlKey(sf.Tag)
		path := elemPathPrefix + "." + y
		doc := manifestYAMLDoc[path]
		if strings.TrimSpace(doc) == "" {
			return fmt.Errorf("manifestYAMLDoc missing %q", path)
		}
		ft := sf.Type
		for ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if fi == 0 {
			writeYAMLCommentBlock(level+1, docWithLogicalType(path, doc, typeByPath), b)
			b.WriteString(padList)
			b.WriteString("- ")
			b.WriteString(y)
			b.WriteString(": ")
			if err := writeYAMLScalarNoNewline(ft, path, b); err != nil {
				return err
			}
			b.WriteString("\n")
			continue
		}
		writeYAMLCommentBlock(level+2, docWithLogicalType(path, doc, typeByPath), b)
		b.WriteString(padCont)
		b.WriteString(y)
		b.WriteString(": ")
		if err := writeYAMLScalarNoNewline(ft, path, b); err != nil {
			return err
		}
		b.WriteString("\n")
	}
	return nil
}

func exportedYAMLFields(t reflect.Type) ([]reflect.StructField, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not struct")
	}
	var out []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		if k := yamlKey(sf.Tag); k != "" && k != "-" {
			out = append(out, sf)
		}
	}
	return out, nil
}

func writeYAMLScalarNoNewline(ft reflect.Type, path string, b *strings.Builder) error {
	switch ft.Kind() {
	case reflect.String:
		fmt.Fprintf(b, "%q", referenceStringValue(path))
	case reflect.Slice:
		el := ft.Elem()
		for el.Kind() == reflect.Ptr {
			el = el.Elem()
		}
		if el.Kind() == reflect.String {
			b.WriteString("[]")
		} else {
			return fmt.Errorf("unsupported slice element for path %s", path)
		}
	default:
		return fmt.Errorf("unsupported scalar kind %s for path %s", ft.Kind(), path)
	}
	return nil
}

func referenceStringValue(path string) string {
	switch path {
	case "version":
		return "v1"
	case "evidence.coverage.inputs.adapter":
		return "hado-json"
	case "evidence.coverage.inputs.path":
		return "coverage-metrics.json"
	default:
		return ""
	}
}

func writeYAMLCommentBlock(level int, doc string, b *strings.Builder) {
	pad := indentSpaces(level)
	for _, line := range strings.Split(doc, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			b.WriteString(pad)
			b.WriteString("#\n")
			continue
		}
		b.WriteString(pad)
		b.WriteString("# ")
		b.WriteString(line)
		b.WriteString("\n")
	}
}

func docWithLogicalType(path, doc string, typeByPath map[string]string) string {
	if typ, ok := typeByPath[path]; ok && typ != "" {
		return doc + " （論理型: " + typ + "）"
	}
	return doc
}

func indentSpaces(level int) string {
	return strings.Repeat("  ", level)
}

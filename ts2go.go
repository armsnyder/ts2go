// Package ts2go is a modular and highly customizable code generator for
// converting TypeScript type definitions into Golang code.
//
// The [Generate] function is the primary entry point for the package. It
// accepts a TypeScript source file as input and writes the generated Golang
// code to an output writer.
//
// The [WithMixin] function can be used to customize the data that is passed to
// the templates. This is useful for adding custom data to the parsed types
// before they are rendered.
//
// The [WithTemplateOverrideDir] function can be used to specify a directory
// that contains template overrides. This is useful for customizing the
// generated code without modifying the built-in templates. The templates used
// by this package are highly modular, allowing you to override only the parts
// that you need.
package ts2go

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"text/template"

	"github.com/armsnyder/typescript-ast-go/parser"
)

//go:embed internal/templates
var templateFS embed.FS

// TemplateData is a wrapper around all of the data which is passed to the
// templates. It is used by [Mixin] functions to mutate the data before it is
// rendered.
type TemplateData struct {
	SkipHeader  bool
	PackageName string
	Structs     []*Struct
	TypeAliases []*TypeAlias
	ConstGroups []*ConstGroup
}

// CustomData contains arbitrary data that can be used by template overrides.
type CustomData map[string]any

// Struct is the data model for the struct.tmpl template.
type Struct struct {
	Name       string
	Doc        []string
	Embeds     []string
	Fields     []*Field
	CustomData CustomData
}

// Field is the data model for a field within a struct.
type Field struct {
	Name       string
	Doc        []string
	Type       string
	IsPointer  bool
	JSONName   string
	OmitEmpty  bool
	CustomData CustomData
}

// TypeAlias is the data model for the type_alias.tmpl template.
type TypeAlias struct {
	Name       string
	Doc        []string
	Type       string
	CustomData CustomData
}

// ConstGroup is the data model for the const_group.tmpl template.
type ConstGroup struct {
	Doc        []string
	CustomData CustomData
}

var DefaultPackageName = "types"

// Generate converts the type definitions in the TypeScript source code into
// Golang and writes the result to the output writer.
func Generate(source io.Reader, output io.Writer, opts ...Option) error {
	g := &generator{
		source: source,
		output: output,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g.generate()
}

type generator struct {
	source             io.Reader
	output             io.Writer
	mixins             []Mixin
	templateOverrideFS fs.FS
}

func (g *generator) generate() error {
	templateData, err := g.parseSource()
	if err != nil {
		return fmt.Errorf("failed to parse source: %w", err)
	}

	g.applyMixins(templateData)

	tmpl, err := g.createTemplate()
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	if err := g.writeOutput(tmpl, templateData); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func (g *generator) parseSource() (*TemplateData, error) {
	sourceBytes, err := io.ReadAll(g.source)
	if err != nil {
		return nil, fmt.Errorf("failed to read source: %w", err)
	}

	sourceFile := parser.Parse(sourceBytes)

	data, err := parseSourceFile(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source file into template data: %w", err)
	}

	return data, nil
}

func (g *generator) applyMixins(templateData *TemplateData) {
	for _, mixin := range g.mixins {
		mixin(templateData)
	}
}

func (g *generator) createTemplate() (*template.Template, error) {
	tmpl, err := template.ParseFS(templateFS, "internal/templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse builtin templates: %w", err)
	}

	// TODO: Implement template overrides
	_ = g.templateOverrideFS

	return tmpl, nil
}

func (g *generator) writeOutput(tmpl *template.Template, templateData *TemplateData) error {
	if err := tmpl.ExecuteTemplate(g.output, "output.tmpl", templateData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

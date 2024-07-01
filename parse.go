package ts2go

import (
	"strings"

	"github.com/armsnyder/typescript-ast-go/ast"
)

func parseSourceFile(sourceFile *ast.SourceFile) (*TemplateData, error) {
	data := &TemplateData{
		PackageName: DefaultPackageName,
	}

	ast.Inspect(sourceFile, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.TypeAliasDeclaration:
			handleTypeAliasDeclaration(data, node)

		case *ast.InterfaceDeclaration:
			handleInterfaceDeclaration(data, node)
		}

		return true
	})

	return data, nil
}

func handleTypeAliasDeclaration(data *TemplateData, node *ast.TypeAliasDeclaration) {
	switch typ := node.Type.(type) {
	case *ast.TypeLiteral:
		data.Structs = append(data.Structs, createStruct(node.Name.Text, node.LeadingComment, typ.Members))

	case *ast.TypeReference:
		switch typeName := typ.TypeName.(type) {
		case *ast.Identifier:
			data.TypeAliases = append(data.TypeAliases, createTypeAlias(node.Name.Text, typeName.Text, node.LeadingComment))

		case *ast.TypeLiteral:
			data.Structs = append(data.Structs, createStruct(node.Name.Text, node.LeadingComment, typeName.Members))
		}
	}
}

func handleInterfaceDeclaration(data *TemplateData, node *ast.InterfaceDeclaration) {
	stru := createStruct(node.Name.Text, node.LeadingComment, node.Members)

	for _, h := range node.HeritageClauses {
		for _, t := range h.Types {
			stru.Embeds = append(stru.Embeds, t.Expression.Text)
		}
	}

	data.Structs = append(data.Structs, stru)
}

func createStruct(name, leadingComment string, members []ast.Signature) *Struct {
	stru := &Struct{
		Name: name,
	}

	if leadingComment != "" {
		stru.Doc = strings.Split(leadingComment, "\n")
	}

	for _, mem := range members {
		switch mem := mem.(type) {
		case *ast.PropertySignature:
			addFieldFromPropertySignature(stru, mem)

		case *ast.IndexSignature:
			// TODO: Handle index signatures
		}
	}

	return stru
}

func addFieldFromPropertySignature(stru *Struct, prop *ast.PropertySignature) {
	field := &Field{
		JSONName:  prop.Name.Text,
		Name:      strings.ToUpper(prop.Name.Text[:1]) + prop.Name.Text[1:],
		OmitEmpty: prop.QuestionToken,
		IsPointer: prop.QuestionToken,
	}

	if prop.LeadingComment != "" {
		field.Doc = strings.Split(prop.LeadingComment, "\n")
	} else if prop.TrailingComment != "" {
		field.Doc = []string{prop.TrailingComment}
	}

	switch t := prop.Type.(type) {
	case *ast.TypeReference:
		switch t := t.TypeName.(type) {
		case *ast.Identifier:
			switch t.Text {
			case "boolean":
				field.Type = "bool"

			default:
				field.Type = t.Text
			}

		default:
			field.Type = "any"
		}

	default:
		field.Type = "any"
	}

	stru.Fields = append(stru.Fields, field)
}

func createTypeAlias(name, referencedName, leadingComment string) *TypeAlias {
	alias := &TypeAlias{
		Name: name,
		Type: referencedName,
	}

	if leadingComment != "" {
		alias.Doc = strings.Split(leadingComment, "\n")
	}

	return alias
}

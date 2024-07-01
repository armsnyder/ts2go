package ts2go_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/armsnyder/ts2go"
)

func Example() {
	source := strings.NewReader(`
/**
 * My type
 */
type Foo = { bar?: string /* My field */ }`)

	output := &bytes.Buffer{}

	err := ts2go.Generate(source, output)

	if err != nil {
		panic(err)
	}

	fmt.Print(output)

	// Output:
	// // Code generated by ts2go. DO NOT EDIT.
	// package types
	//
	// // My type
	// type Foo struct {
	//	// My field
	//	Bar *string `json:"bar,omitempty"`
	// }
}

func Example_mixins() {
	source := strings.NewReader(`
/**
 * My type
 */
type Foo = { bar?: string /* My field */ }`)

	output := &bytes.Buffer{}

	err := ts2go.Generate(source, output, ts2go.WithMixin(
		ts2go.SetPackageName("mypackage"),
		ts2go.SkipOptionalPointer(),
		func(td *ts2go.TemplateData) {
			td.Structs[0].Name = "Override"
		},
	))

	if err != nil {
		panic(err)
	}

	fmt.Print(output)

	// Output:
	// // Code generated by ts2go. DO NOT EDIT.
	// package mypackage
	//
	// // My type
	// type Override struct {
	//	// My field
	//	Bar string `json:"bar,omitempty"`
	// }
}
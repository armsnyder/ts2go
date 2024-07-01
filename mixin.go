package ts2go

// Mixin is used to modify or add custom data to the [TemplateData] before it
// is rendered.
type Mixin func(*TemplateData)

// SkipOptionalPointer specifies that optional fields should not be
// represented as pointers in the generated Go code.
func SkipOptionalPointer() Mixin {
	return func(data *TemplateData) {
		for _, s := range data.Structs {
			for _, f := range s.Fields {
				f.IsPointer = false
			}
		}
	}
}

// SkipHeader specifies that the generated Go code should not include the
// header comment and package declaration. This is useful if you want to run
// [Generate] more than once and concatenate the results.
func SkipHeader() Mixin {
	return func(data *TemplateData) {
		data.SkipHeader = true
	}
}

// SetPackageName specifies the package name that should be used in the
// generated Golang code. If not specified, [DefaultPackageName] is used.
func SetPackageName(name string) Mixin {
	return func(data *TemplateData) {
		data.PackageName = name
	}
}

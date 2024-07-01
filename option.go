package ts2go

import (
	"io/fs"
	"os"
)

// Option customizes the behavior of [Generate].
type Option func(*generator)

// WithMixin specifies one or more [Mixin]s to customize the data that is
// passed to the templates.
func WithMixin(mixins ...Mixin) Option {
	return func(g *generator) {
		g.mixins = append(g.mixins, mixins...)
	}
}

// WithTemplateOverrideDir specifies a directory which can contain template
// overrides.
//
// You may override as many or as few templates as you like. This is useful in
// conjunction with mixins. You can add custom data to the parsed types, and
// then reference that custom data inside your template override. See the
// templates directory for the list of templates that can be overridden.
func WithTemplateOverrideDir(dir string) Option {
	return WithTemplateOverrideFS(os.DirFS(dir))
}

// WithTemplateOverrideFS is the same as [WithTemplateOverrideDir], but takes a
// file system interface instead of a directory.
func WithTemplateOverrideFS(fsys fs.FS) Option {
	return func(g *generator) {
		g.templateOverrideFS = fsys
	}
}

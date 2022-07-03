package template

import (
	t "html/template"

	"github.com/maxgio92/krawler/pkg/matrix"
)

type TemplatePart struct {
	matrix.Column
	TemplateString  string
	MatchedVariable string
	Template        *t.Template
}

package template

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	t "html/template"

	"github.com/maxgio92/krawler/pkg/matrix"
)

// Returns a list of strings of executed templates from a template string
// input, by applying an arbitrary variables inventory with multiple values.
// The expecetd arguments are:

// - templateString: the input string of the template to execute.
// - inventory: the inventory map of the variable data structure to apply the
// template to. The key of the map is the name of the variable, that should
// match related annotation in the template. Each map item is a slice where in each slice item is
// a single variable value.
//
// The result is multiple templates from a single template string and multiple
// arbitrary variable values.
func MultiplexAndExecute(templateString string, inventory map[string][]interface{}) ([]string, error) {
	supportedVariables, err := getSupportedVariables(templateString)
	if err != nil {
		return nil, err
	}

	templateRegex, err := generateTemplateRegex(supportedVariables)
	if err != nil {
		return nil, err
	}
	templatePattern := regexp.MustCompile(templateRegex)

	ss, err := cutTemplateString(templateString, closeDelimiter)
	if err != nil {
		return nil, err
	}

	templateParts := []TemplatePart{}

	for _, s := range ss {

		// match are the template parts matched against the template regex.
		templatePartMatches := templatePattern.FindStringSubmatch(s)

		// name is the variable data structure to apply the template part to.
		for i, variableName := range templatePattern.SubexpNames() {

			// discard first variable name match and ensure a template part matched.
			if i > 0 && i <= len(templatePartMatches) && templatePartMatches[i] != "" {
				y := len(templateParts)

				templateParts = append(templateParts, TemplatePart{
					TemplateString:  templatePartMatches[i],
					MatchedVariable: variableName,
				})

				templateParts[y].Points = []string{}
				templateParts[y].TemplateString = strings.ReplaceAll(
					templateParts[y].TemplateString,
					openDelimiter+` `+cursor+variableName+` `+closeDelimiter,
					openDelimiter+` `+cursor+` `+closeDelimiter,
				)
				templateParts[y].Template = t.New(fmt.Sprintf("%d", y))
				templateParts[y].Template, err = templateParts[y].Template.Parse(templateParts[y].TemplateString)
				if err != nil {
					return nil, err
				}

				// for each item (variable name) of MatchedVariable
				// compose one Template and `execute()` it
				for _, value := range inventory[variableName] {
					o := new(bytes.Buffer)
					err = templateParts[y].Template.Execute(o, value)
					if err != nil {
						return nil, err
					}

					templateParts[y].Points = append(templateParts[y].Points.([]string), o.String())
				}
			}
		}
	}

	matrixColumns := []matrix.Column{}

	for _, part := range templateParts {
		matrixColumns = append(matrixColumns, part.Column)
	}

	if len(matrixColumns) <= 0 {
		return nil, fmt.Errorf("cannot multiplex template: the template contains syntax errors")
	}

	result, err := matrix.GetColumnOrderedCombinationRows(matrixColumns)
	if err != nil {
		return nil, err
	}

	return result, nil
}

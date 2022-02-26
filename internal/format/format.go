package format

import (
	"bufio"
	"encoding/json"

	"gopkg.in/yaml.v2"

	"github.com/olekukonko/tablewriter"
)

type Type string

const (
	Text Type = "text"
	JSON Type = "json"
	YAML Type = "yaml"
)

func Encode(output *bufio.Writer,
	objects interface{},
	format Type) (*bufio.Writer, error) {
	switch format {
	case JSON:
		return encodeJSON(output, objects)
	case Text:
		return encodeText(output, objects)
	case YAML:
		return encodeYAML(output, objects)
	default:
		return encodeText(output, objects)
	}
}

func encodeJSON(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	json, err := json.Marshal(objects)
	if err != nil {
		return nil, err
	}

	_, err = output.Write(json)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func encodeYAML(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	yaml, err := yaml.Marshal(objects)
	if err != nil {
		return nil, err
	}

	_, err = output.Write(yaml)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func encodeText(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	return encodeTableFromStructs(output, objects)
}

func encodeTableFromStructs(output *bufio.Writer, objects interface{}) (*bufio.Writer, error) {
	printer := tablewriter.NewWriter(output)

	printer.SetStructs(objects)
	printer.Render()

	return output, nil
}

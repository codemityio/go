package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/codemityio/go/pkg/fsm"
	"github.com/codemityio/go/pkg/fsm/dsl"
)

//go:embed diagram.tpl
var diagram string

// DefaultGenerator a fsm diagram generator struct.
type DefaultGenerator struct {
	start, stop, state, link string
}

func (d *DefaultGenerator) Generate(workflow *dsl.Workflow) ([]byte, error) {
	smw := fsm.New(workflow)

	// sort states and edges
	sort.Slice(workflow.States, func(i, j int) bool {
		return workflow.States[i].Name < workflow.States[j].Name
	})

	sort.Slice(workflow.Edges, func(i, j int) bool {
		return workflow.Edges[i].From < workflow.Edges[j].From
	})

	state := smw.GetInitialState()
	if state == nil {
		return nil, fmt.Errorf("%w: initial state not found", ErrPkg)
	}

	params := map[string]any{
		"Initial": state.Name,
		"Final":   smw.GetFinalStates(),
		"States":  workflow.States,
		"Edges":   workflow.Edges,
		"Start":   d.start,
		"Stop":    d.stop,
		"State":   d.state,
		"Link":    d.link,
	}

	tpl := template.New("diagram").Funcs(template.FuncMap{
		"trimHash": func(value string) string {
			return strings.Trim(value, "#")
		},
	})

	tpl, err := tpl.Parse(diagram)
	if err != nil {
		return nil, fmt.Errorf("%w: error parsing a template: %w", ErrPkg, err)
	}

	buff := bytes.Buffer{}

	err = tpl.Execute(&buff, params)
	if err != nil {
		return nil, fmt.Errorf("%w: error while generating puml diagram: %w", ErrPkg, err)
	}

	return buff.Bytes(), nil
}

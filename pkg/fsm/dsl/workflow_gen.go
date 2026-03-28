package dsl

type Workflow struct {
	Edges  []Edge  `json:"edges"`
	States []State `json:"states"`
}

type Edge struct {
	Colour      string `json:"colour"`
	Description string `json:"description"`
	From        string `json:"from"`
	To          string `json:"to"`
}

type State struct {
	Colour      string `json:"colour"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

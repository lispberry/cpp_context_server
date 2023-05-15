package semantic

type Ast interface{}

type File struct {
	Functions []*Function `json:"functions"`
}

type Function struct {
	Name       string      `json:"name"`
	Arguments  []*Argument `json:"arguments"`
	LineNumber int         `json:"line_number"`
}

type Argument struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

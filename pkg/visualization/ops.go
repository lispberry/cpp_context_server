package visualization

import "encoding/json"

type RawOp struct {
	Kind string `json:"kind"`
	Data json.RawMessage
}

const (
	NewListPointerKind = "NewListPointer"
)

type Op interface{}

type NewListPointer struct {
	Name string
}

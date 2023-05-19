package visualization

import "encoding/json"

type RawOp struct {
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

const (
	NewListPointerKind         = "NewListPointer"
	NewListNodeKind            = "NewListNode"
	SetListNodeNextKind        = "SetListNodeNext"
	SetListNodeValueKind       = "SetListNodeValue"
	SetListPointerValueKind    = "SetListPointerValue"
	DereferenceListPointerKind = "DereferenceListPointer"
)

type Op interface{}

type NewListPointer struct {
	Name string `json:"name"`
}

type SetListPointerValue struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type NewListNode struct {
	Address string `json:"address"`
}

type SetListNodeNext struct {
	Address string `json:"address"`
	Next    string `json:"next"`
}

type SetListNodeValue struct {
	Address string `json:"address"`
	Value   string `json:"value"`
}

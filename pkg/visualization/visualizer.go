package visualization

import (
	"encoding/json"
	"github.com/lispberry/viz-service/pkg/dot"
)

type Visualizer struct {
	pointers map[dot.Ref]string
	memGraph *dot.MemoryGraph
}

func NewVisualizer() (*Visualizer, error) {
	return &Visualizer{
		memGraph: dot.NewMemoryGraph(),
	}, nil
}

func (v *Visualizer) Apply(ops []RawOp) ([]string, error) {
	v.memGraph = dot.NewMemoryGraph()

	for _, rawOp := range ops {
		op, err := v.newOp(rawOp)
		if err != nil {
			return nil, err
		}
		v.applyOp(op)
	}

	return v.memGraph.Changes(), nil
}

func (v *Visualizer) newOp(op RawOp) (Op, error) {
	var action interface{}
	switch op.Kind {
	case NewListPointerKind:
		action = new(NewListPointer)
	case SetListPointerValueKind:
		action = new(SetListPointerValue)
	case NewListNodeKind:
		action = new(NewListNode)
	case SetListNodeNextKind:
		action = new(SetListNodeNext)
	case SetListNodeValueKind:
		action = new(SetListNodeValue)
	}

	err := json.Unmarshal(op.Data, action)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (v *Visualizer) applyOp(op Op) {
	switch opd := op.(type) {
	case *NewListPointer:
		v.memGraph.Pointer(opd.Name, dot.NullptrRef)
	case *SetListPointerValue:
		ptr := v.memGraph.GetPointer(opd.Name)
		ptr.SetAddress(dot.Ref(opd.Address))
	case *SetListNodeNext:
		node := v.memGraph.GetListNode(dot.Ref(opd.Address))
		node.SetNext(dot.Ref(opd.Next))
	case *SetListNodeValue:
		node := v.memGraph.GetListNode(dot.Ref(opd.Address))
		node.SetData(opd.Value)
	case *NewListNode:
		v.memGraph.ListNode(dot.Ref(opd.Address), "0", dot.NullptrRef)
	default:
		println("here")
	}
}

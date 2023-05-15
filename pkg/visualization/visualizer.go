package visualization

import (
	"encoding/json"
	"github.com/lispberry/viz-service/pkg/dot"
)

type Visualizer struct {
	pointers map[dot.Ref]string
	memGraph *dot.MemoryGraph
}

func NewVisualizer() *Visualizer {
	return &Visualizer{
		memGraph: dot.NewMemoryGraph(),
	}
}

func (v *Visualizer) Apply(ops []RawOp) ([]string, error) {
	var dots []string

	for _, rawOp := range ops {
		op, err := v.newOp(rawOp)
		if err != nil {
			return nil, err
		}

		dots = append(dots, v.applyOp(op))
	}

	return dots, nil
}

func (v *Visualizer) newOp(op RawOp) (Op, error) {
	var action interface{}
	switch op.Kind {
	case NewListPointerKind:
		action = new(NewListPointer)
	}

	err := json.Unmarshal(op.Data, action)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (v *Visualizer) applyOp(op Op) string {
	switch opd := op.(type) {
	case NewListPointer:
		return opd.Name
	default:
		return ""
	}
}

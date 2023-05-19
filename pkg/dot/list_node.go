package dot

type ListNode struct {
	mem     *MemoryGraph
	node    *Node
	Address Ref
	Data    string
	Next    Ref
}

const (
	addressPort = "address"
	dataPort    = "data"
	nextPort    = "next"
)

func defaultDescAttrs() map[string]string {
	return map[string]string{
		"border": "0",
	}
}

func NewListNode(mem *MemoryGraph, address Ref, data string, next Ref) *ListNode {
	node := NewNode(address, 2, 3)
	node.Attrs = map[string]string{
		"border":      "0",
		"cellspacing": "0",
		"cellborder":  "1",
	}

	// TODO(ivo): Ugly
	nextAttrs := defaultAttrs(nextPort)
	nextValue := next.Value()
	if next == NullptrRef {
		nextAttrs["bgColor"] = nullptrColor
		nextValue = ""
	}

	node.Table[0] = []Cell{
		Cell{
			Value: address.Value(),
			Attrs: defaultAttrs(addressPort),
		},
		Cell{
			Value: data,
			Attrs: defaultAttrs(dataPort),
		},
		Cell{
			Value: nextValue,
			Attrs: nextAttrs,
		},
	}
	node.Table[1] = []Cell{
		Cell{
			Value: "addr",
			Attrs: defaultDescAttrs(),
		},
		Cell{
			Value: "data",
			Attrs: defaultDescAttrs(),
		},
		Cell{
			Value: "next",
			Attrs: defaultDescAttrs(),
		},
	}

	return &ListNode{
		mem:     mem,
		node:    node,
		Address: address,
		Data:    data,
		Next:    next,
	}
}

func (n *ListNode) SetData(data string) {
	n.node.Table[0][1].Attrs["bgColor"] = defaultChangeColor
	n.mem.changed()

	n.Data = data
	n.node.Table[0][1].Value = data
	n.mem.changed()

	delete(n.node.Table[0][1].Attrs, "bgColor")
	n.mem.changed()
}

func (n *ListNode) SetNext(next Ref) {
	if next == NullptrRef {
		n.Next = "0x0"
		n.node.Table[0][2].Attrs["bgColor"] = nullptrColor
		n.node.Table[0][2].Value = ""
		n.mem.changed()
		return
	}

	n.node.Table[0][2].Attrs["bgColor"] = defaultChangeColor
	n.mem.changed()

	n.Next = next
	n.node.Table[0][2].Value = next.Value()
	n.mem.changed()

	delete(n.node.Table[0][2].Attrs, "bgColor")
	n.mem.changed()
}

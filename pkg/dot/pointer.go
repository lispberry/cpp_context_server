package dot

type Pointer struct {
	mem     *MemoryGraph
	node    *Node
	Name    string
	Address Ref
}

const (
	pointerValuePort = "value"
	pointerNamePort  = "name"
)

func newPointer(mem *MemoryGraph, name string, address Ref) *Pointer {
	node := NewNode(Ref(name), 1, 2)
	node.Attrs = map[string]string{
		"border":      "0",
		"cellspacing": "0",
		"cellborder":  "1",
	}

	// TODO(ivo): Ugly
	addrAttrs := defaultAttrs(pointerValuePort)
	addrValue := address.Value()
	if address == NullptrRef {
		addrAttrs["bgColor"] = nullptrColor
		addrValue = ""
	}

	node.Table[0] = []Cell{
		Cell{
			Value: name,
			Attrs: defaultAttrs(pointerNamePort),
		},
		Cell{
			Value: addrValue,
			Attrs: addrAttrs,
		},
	}

	return &Pointer{
		mem:     mem,
		node:    node,
		Name:    name,
		Address: address,
	}
}

func (p *Pointer) SetAddress(address Ref) {
	if address == NullptrRef {
		p.Address = "0x0"
		p.node.Table[0][1].Attrs["bgColor"] = nullptrColor
		p.node.Table[0][1].Value = ""
		p.mem.changed()
		return
	}

	p.node.Table[0][1].Attrs["bgColor"] = defaultChangeColor
	p.mem.changed()

	p.Address = address
	p.node.Table[0][1].Value = address.Value()
	p.mem.changed()

	delete(p.node.Table[0][1].Attrs, "bgColor")
	p.mem.changed()
}

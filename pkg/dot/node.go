package dot

import (
	"fmt"
	"strings"
)

type Cell struct {
	Attrs map[string]string
	Value string
}

func (cell *Cell) String() string {
	var attrs strings.Builder
	for k, v := range cell.Attrs {
		attrs.WriteString(fmt.Sprintf(`%s="%s" `, k, v))
	}

	return fmt.Sprintf("<td %s>%s</td>", attrs.String(), cell.Value)
}

type Node struct {
	Ref   Ref
	Attrs map[string]string
	Table [][]Cell
}

const nullptrColor = "##B8B3B2"
const defaultChangeColor = "#F9FF67"

func defaultAttrs(port string) map[string]string {
	return map[string]string{
		"port":   port,
		"width":  "28",
		"height": "36",
	}
}

func NewNode(ref Ref, rows int, columns int) *Node {
	var table = make([][]Cell, rows)
	for i := range table {
		table[i] = make([]Cell, columns)
	}

	return &Node{
		Ref:   ref,
		Attrs: map[string]string{},
		Table: table,
	}
}

func (n *Node) String() string {
	var tableAttrs strings.Builder
	for k, v := range n.Attrs {
		tableAttrs.WriteString(fmt.Sprintf(`%s="%s" `, k, v))
	}

	var table strings.Builder
	table.WriteString(fmt.Sprintf("<<table %s>", tableAttrs.String()))
	for _, row := range n.Table {
		table.WriteString("<tr>")
		for _, cell := range row {
			table.WriteString(cell.String())
		}
		table.WriteString("</tr>")
	}
	table.WriteString("</table>>")

	return table.String()
}

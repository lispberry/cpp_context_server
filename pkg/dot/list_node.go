package dot

import "fmt"

type ListNode struct {
	Address string
	Data    string
	Next    string
}

const addressPort = "ref1"
const dataPort = "data"
const nextPort = "ref2"

func (node *ListNode) Table() string {
	const table = `<<table border="0" cellspacing="0" cellborder="1">
		<tr>
			<td port="ref1" width="28" height="36">%s</td>
			<td port="data" width="28" height="36">%s</td>
			<td port="ref2" width="28" height="36">%s</td>
		</tr>
		<tr>
			<td BORDER="0">addr</td>
			<td BORDER="0">val</td>
			<td BORDER="0">next</td>
		</tr>
	</table>>`

	return fmt.Sprintf(table, node.Address, node.Data, node.Next)
}

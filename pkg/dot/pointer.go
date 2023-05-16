package dot

import "fmt"

type Pointer struct {
	Name    string
	Address Ref
}

func (p *Pointer) Table() string {
	var address string
	if p.Address == "0x0" || p.Address == "" {
		address = `<td port="address" width="28" height="36" bgcolor="#C1FF83"></td>`
	} else {
		address = fmt.Sprintf(`<td port="address" width="28" height="36" bgcolor="#C1FF83">%s</td>`, p.Address.Value())
	}

	const table = `<<table border="0" cellspacing="0" cellborder="1">
		<tr>
			<td port="name" width="28" height="36">%s</td>
			%s
		</tr>
	</table>>`

	return fmt.Sprintf(table, p.Name, address)
}

func (p *Pointer) SetAddress(address string) {

}

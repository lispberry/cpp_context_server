package visualization

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestVisualizer(t *testing.T) {
	vis, _ := NewVisualizer()

	vis.Apply([]RawOp{
		{
			Kind: NewListPointerKind,
			Data: []byte(`{"name": "Hello"}`),
		},
		{
			Kind: NewListNodeKind,
			Data: []byte(`{"address": "0xFFFF"}`),
		},
		{
			Kind: NewListNodeKind,
			Data: []byte(`{"address": "0xFFFF1"}`),
		},
		{
			Kind: SetListPointerValueKind,
			Data: []byte(`{"name":"Hello","address": "0xFFFF1"}`),
		},
		{
			Kind: SetListNodeNextKind,
			Data: []byte(`{"address": "0xFFFF", "next": "0xFFFF1"}`),
		},
	})

	test := `
[
{"data":{"name":"head"},"kind":"NewListPointer"},
{"data":{"address":"0x6000036a8030"},"kind":"NewListNode"},
{"data":{"name":"scd"},"kind":"NewListPointer"},
{"data":{"address":"0x6000036a8040"},"kind":"NewListNode"},
{"data":{"address":"0x6000036a8030","next":"0x6000036a8040"},"kind":"SetListNodeNext"},
{"data":{"name":"thrd"},"kind":"NewListPointer"},
{"data":{"address":"0x6000036a8050"},"kind":"NewListNode"},
{"data":{"address":"0x6000036a8040","next":"0x6000036a8050"},"kind":"SetListNodeNext"},
{"data":{"address":"0x6000036a8040","name":"head"},"kind":"SetListPointerValue"},
{"data":{"address":"0x6000036a8050","name":"head"},"kind":"SetListPointerValue"},
{"data":{"name":"var"},"kind":"NewListPointer"},
{"data":{"address":"0x6000036a8060"},"kind":"NewListNode"},{"data":{"address":"0x6000036a8060","name":"var"},"kind":"SetListPointerValue"},
{"data":{"address":"0x6000036a8060","value":"123"},"kind":"SetListNodeValue"},
{"data":{"address":"0x6000036a8050","next":"0x6000036a8060"},"kind":"SetListNodeNext"}
]
	`
	//,
	//{"data":{"address":"0x600002eac030","next":"0x0"},"kind":"SetListNodeNext"},
	//{"data":{"address":"0x600002eac030","value":"1234"},"kind":"SetListNodeValue"}]`
	var ops []RawOp
	json.Unmarshal([]byte(test), &ops)

	res, _ := vis.Apply(ops)

	fmt.Println("[")
	for _, dot := range res {
		fmt.Printf("`%s`,", dot)
	}
	fmt.Println("]")

	//fmt.Printf("%v", res)
}

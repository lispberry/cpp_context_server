package visualization

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestVisualizer(t *testing.T) {
	vis, _ := NewVisualizer()

	test := `
[
{"data":{"name":"head"},"kind":"NewListPointer"},
{"data":{"address":"0x6000036a8030"},"kind":"NewListNode"},
{"data":{"address":"0x6000036a8040"},"kind":"NewListNode"},
{"data":{"address":"0x6000036a8030","next":"0x6000036a8040"},"kind":"SetListNodeNext"},
{"data":{"address":"0x6000036a8050"},"kind":"NewListNode"},
{"data":{"address":"0x6000036a8040","next":"0x6000036a8050"},"kind":"SetListNodeNext"},
{"data":{"address":"0x6000036a8040","name":"head"},"kind":"SetListPointerValue"},
{"data":{"address":"0x6000036a8050","name":"head"},"kind":"SetListPointerValue"},
{"data":{"name":"var"},"kind":"NewListPointer"},
{"data":{"address":"0x6000036a8060"},"kind":"NewListNode"},{"data":{"address":"0x6000036a8060","name":"var"},"kind":"SetListPointerValue"},
{"data":{"address":"0x6000036a8060","value":"123"},"kind":"SetListNodeValue"},
{"data":{"address":"0x6000036a8050","next":"0x6000036a8060"},"kind":"SetListNodeNext"}
]`
	//
	//	test1 := `[
	//{"data":{"address":"0xaaaaaaaefeb0"},"kind":"NewListNode"},
	//{"data":{"name":"head"},"kind":"NewListPointer"},{"data":{"address":"0xaaaaaaaefeb0","name":"head"},"kind":"SetListPointerValue"},
	//{"data":{"name":"var"},"kind":"NewListPointer"},
	//{"data":{"address":"0xaaaaaaaf0290"},"kind":"NewListNode"},{"data":{"address":"0xaaaaaaaf0290","name":"var"},"kind":"SetListPointerValue"},
	//{"data":{"address":"0xaaaaaaaf0290","value":"10"},"kind":"SetListNodeValue"},
	//{"data":{"address":"0xaaaaaaaefeb0","next":"0xaaaaaaaf0290"},"kind":"SetListNodeNext"}
	//]
	//`
	//,
	//{"data":{"address":"0x600002eac030","next":"0x0"},"kind":"SetListNodeNext"},
	//{"data":{"address":"0x600002eac030","value":"1234"},"kind":"SetListNodeValue"}]`
	var ops []RawOp
	err := json.Unmarshal([]byte(test), &ops)
	if err != nil {
		t.Fail()
	}

	err = vis.Apply(ops)
	if err != nil {
		t.Fail()
	}

	fmt.Println("[")
	for _, dot := range vis.Changes() {
		fmt.Printf("`%s`,", dot)
	}
	fmt.Println("]")

	//fmt.Printf("%v", res)
}

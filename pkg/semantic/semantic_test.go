package semantic

import (
	"encoding/json"
	"fmt"
	"github.com/lispberry/viz-service/pkg/evaluation"
	"github.com/lispberry/viz-service/pkg/visualization"
	"testing"
)

func TestProgramReflect(t *testing.T) {
	program := NewProgram(`
		struct List
		{
			int x;
			int y;
		};
		
		int program(List * arg1, int arg2) {}

		int main()
		{
			return 0;
		}
	`)
	file, err := program.Reflect()
	if err != nil {
		t.Fail()
	}

	if file.Functions[0].Name != "program" {
		t.Fail()
	}
	if file.Functions[1].Name != "main" {
		t.Fail()
	}
}

/*
  Pointer<List> h("head");     // List * head;
  SetListPointerValue(h);      // head = new List;
  Pointer<List> w("wow");      // List * wow;
  SetListPointerValue(w);      // wow = new List;
  SetListNodeNext(h, w);       // head->next = w;
  SetListNodeValue(h, 1234);   // head->value = 1234;
*/

func TestInstrumentProgram(t *testing.T) {
	fmt.Println(InstrumentProgram("List * test = new List;"))
	fmt.Println(InstrumentProgram("void Program(List *x, List *y)"))
	fmt.Println(InstrumentProgram(`
void insertEnd(List * head, int val)
{
  while (head->n != nullptr)
  {
    head = head->n;
  }
  List * var;
  var = new List;
  var->val = val;
  head->n = var;
}
	`))
}

func TestExecutableProgram(t *testing.T) {
	const sourceCode = `
void insertEnd(List * head, int val)
{
  while (head->n != nullptr)
  {
    head = head->n;
  }
  List * var;
  var = new List;
  var->val = val;
  head->n = var;
}`

	program := NewProgram(sourceCode)
	path, err := program.Executable()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if path == "" {
		t.Fail()
	}

	file, err := program.Reflect()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if file.Functions[0].Name != "insertEnd" {
		t.Fatalf("Expected insertEnd name")
	}
}

func TestNewProgram(t *testing.T) {
	const sourceCode = `
void insertEnd(List * head, int val)
{
  while (head->n != nullptr)
  {
    head = head->n;
  }
  List * var;
  var = new List;
  var->val = val;
  head->n = var;
}`
	program := NewProgram(sourceCode)
	path, err := program.Executable()
	if err != nil {
		t.Fatalf("%v", err)
	}
	if path == "" {
		t.Fail()
	}

	file, err := program.Reflect()
	if err != nil {
		t.Fatalf("%v", err)
	}

	fn := file.Functions[0]

	lev, err := evaluation.New()
	defer lev.Close()
	if err != nil {
		t.Fatalf("%v", err)
	}
	err = lev.EvaluateFile(path)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var args []interface{}
	for _, arg := range fn.Arguments {
		if arg.Type == "int *" {
			args = append(args, fmt.Sprintf(`unmove(CreatePointer("%s"))`, arg.Name))
		} else {
			args = append(args, 10)
		}
	}

	next, err := lev.EvaluateFunc(`i`, args...)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// vis, err := visualization.NewVisualizer()
	if err != nil {
		t.Fatalf("%v", err)
	}

	output := lev.Output()
	for next() {
		fmt.Println(lev.CurrentLine())
		select {
		case line, ok := <-output:
			if ok {
				var ops []visualization.RawOp
				err = json.Unmarshal([]byte(line), &ops)
				if err != nil {
					t.Fatalf("%v", err)
				}
				fmt.Println(line)
			}
		default:
		}
	}
}

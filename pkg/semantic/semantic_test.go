package semantic

import (
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

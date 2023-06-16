package evaluation

import (
	"reflect"
	"testing"
)

func TestEvaluateEmptyProgram(t *testing.T) {
	lev, err := New()
	if err != nil {
		t.Fail()
	}
	defer lev.Close()

	err = lev.EvaluateProgram(`int main() {
		return 0;
	}`)
	if err != nil {
		t.Fail()
	}
}

func TestReadOutput(t *testing.T) {
	lev, err := New()
	if err != nil {
		t.Fail()
	}
	defer lev.Close()

	err = lev.EvaluateProgram(`
	#include <cstdio>	

	int main() {
		return 0;
	}`)
	if err != nil {
		t.Fail()
	}

	_, err = lev.EvaluateExp("puts(\"123\")")
	if err != nil {
		t.Fail()
	}

	select {
	case line, ok := <-lev.Output():
		if !ok || line != "123" {
			t.Fail()
		}
	default:
		t.Fail()
	}
}

func TestEvaluateExp(t *testing.T) {
	lev, err := New()
	if err != nil {
		t.Fail()
	}
	defer lev.Close()

	err = lev.EvaluateProgram(`
		#include <cstdio>

		struct Test {
			int x;
			int y;
		};

		void program()
		{
			Test test{10, 20};
			int *x = new int;
		}

		int main() {
			Test test{10, 20};
			int *x = new int;

		return 0;
		}`)
	if err != nil {
		t.Fail()
	}

	tests := []struct {
		name     string
		exp      string
		expected interface{}
	}{
		{"Add two numbers", "10 + 10", "20"},
		{"Create a struct", "test", "{x = 10, y = 20}"},
	}

	next, err := lev.EvaluateFunc("program", "void(*)()")
	if err != nil {
		t.Fail()
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			next()
			result, err := lev.EvaluateExp(test.exp)
			if err != nil {
				t.Fail()
			}
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("EvaluateExp(%s) == %v; expected %v", test.exp, result, test.expected)
			}
		})
	}
}

func TestEvaluateFunction(t *testing.T) {
	lev, err := New()
	if err != nil {
		t.Fail()
	}
	defer lev.Close()

	err = lev.EvaluateProgram(`
		#include <cstdio>

		int program(int z, int y, int ww)
		{
			int x = 123;
			return 10;
		}

		int main() {
			return 0;
		}
	`)
	if err != nil {
		t.Fail()
	}
	next, err := lev.EvaluateFunc("program", "int(*)(int, int, int)", 10, 20, 30)
	if err != nil {
		t.Errorf("%v", err)
	}
	if line, err := lev.CurrentLine(); err != nil || line != 6 {
		t.Errorf("Expected line %d got %d", 6, line)
	}

	for next() {
	}
}

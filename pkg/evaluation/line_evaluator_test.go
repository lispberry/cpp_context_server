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

func TestEvaluateExp(t *testing.T) {
	lev, err := New()
	if err != nil {
		t.Fail()
	}
	defer lev.Close()

	err = lev.EvaluateProgram(`
		struct Test {
			int x;
			int y;
		};


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
		// TODO(ivo): Add struct values parsing
		{"Create a struct", "test", "{x = 10, y = 20}"},
		{"", "x", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

			printf("hello world\n");
			return 10;
		}

		int main() {
	
			return 0;
		}
	`)
	if err != nil {
		t.Fail()
	}
	next, err := lev.EvaluateFunc("program", 10, 20, 30)
	if err != nil {
		t.Errorf("%v", err)
	}
	if line, err := lev.CurrentLine(); err != nil || line != 5 {
		t.Errorf("Expected line %d got %d", 5, line)
	}
	for next() {
	}
}

func TestEvaluateFunction1(t *testing.T) {

}

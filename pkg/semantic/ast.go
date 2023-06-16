package semantic

import (
	"fmt"
	"strings"
)

type Ast interface{}

type File struct {
	Functions []*Function `json:"functions"`
}

type Function struct {
	Name       string      `json:"name"`
	Arguments  []*Argument `json:"arguments"`
	LineNumber int         `json:"line_number"`
}

type Argument struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// TODO(ivo): FIXME
const listType = "int *"

func (f *Function) Signature() string {
	var args []string
	for _, arg := range f.Arguments {
		if arg.Type == listType {
			args = append(args, "Pointer<List> &")
		} else {
			args = append(args, arg.Type)
		}
	}

	// TODO(ivo): Return Type
	return fmt.Sprintf("void(*)(%s)", strings.Join(args, ", "))
}

// TODO(ivo): Ugly
func (f *Function) Callable(vargs ...interface{}) []interface{} {
	var args []interface{}
	i := 0
	for _, arg := range f.Arguments {
		if arg.Type == listType {
			args = append(args, fmt.Sprintf(`unmove(CreatePointer("%s"))`, arg.Name))
		} else {
			args = append(args, vargs[i])
			i++
		}
	}

	return args
}

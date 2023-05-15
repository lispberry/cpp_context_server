package evaluation

import (
	"fmt"
	"github.com/cyrus-and/gdb"
	"log"
	"strconv"
	"strings"
)

type LineEvaluator struct {
	gdb *gdb.Gdb
}

func New() (*LineEvaluator, error) {
	gdb, err := gdb.New(nil)
	if err != nil {
		return nil, err
	}

	return &LineEvaluator{
		gdb: gdb,
	}, nil
}

func mapGoToCpp(val interface{}) string {
	switch val.(type) {
	case int:
		return strconv.Itoa(val.(int))
	case string:
		return "\"" + val.(string) + "\""
	default:
		panic("couldn't mapGoToCpp")
	}
}

func (lev *LineEvaluator) CurrentFuncName() (string, error) {
	res, err := lev.gdb.Send("stack-info-frame")
	if err != nil {
		return "", err
	}
	payload := res["payload"].(map[string]interface{})
	frame := payload["frame"].(map[string]interface{})
	return frame["func"].(string), nil
}

func (lev *LineEvaluator) CurrentLine() (int, error) {
	res, err := lev.gdb.Send("stack-info-frame")
	if err != nil {
		return 0, err
	}
	payload := res["payload"].(map[string]interface{})
	frame := payload["frame"].(map[string]interface{})
	line := frame["line"].(string)
	return strconv.Atoi(line)
}

func (lev *LineEvaluator) EvaluateProgram(sourceCode string) error {
	program, err := compileProgram(sourceCode)
	if err != nil {
		return err
	}

	_, err = lev.gdb.Send("file-exec-and-symbols", program)
	if err != nil {
		return err
	}

	_, err = lev.gdb.Send("break-insert", "main")
	if err != nil {
		return err
	}

	_, err = lev.gdb.Send("exec-run")
	if err != nil {
		return err
	}

	return nil
}

type ExpResult interface {
}

func (lev *LineEvaluator) EvaluateExp(exp string) (interface{}, error) {
	res, err := lev.gdb.Send("data-evaluate-expression", exp)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)

	payload := res["payload"].(map[string]interface{})
	return payload["value"], nil
}

func (lev *LineEvaluator) EvaluateFunc(funcName string, values ...interface{}) (func() bool, error) {
	mappedValues := make([]string, len(values))
	for i, val := range values {
		mappedValues[i] = mapGoToCpp(val)
	}
	functionCall := funcName + "(" + strings.Join(mappedValues, ", ") + ")"

	_, err := lev.gdb.Send("break-insert", funcName)
	if err != nil {
		return nil, err
	}

	_, err = lev.gdb.Send("data-evaluate-expression", functionCall)
	if err != nil {
		return nil, err
	}

	return func() bool {
		currentFunc, err := lev.CurrentFuncName()
		if err != nil {
			log.Println(err)
			return false
		}
		if currentFunc != funcName {
			return false
		}

		_, err = lev.gdb.Send("exec-next")
		if err != nil {
			log.Println(err)
		}
		return true
	}, nil
}

func (lev *LineEvaluator) Close() error {
	return lev.gdb.Exit()
}

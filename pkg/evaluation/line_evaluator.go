package evaluation

import (
	"fmt"
	"github.com/cyrus-and/gdb"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// TODO(ivo): Remove
func compileProgram(program string) (string, error) {
	process, err := os.CreateTemp("", "program-*")
	if err != nil {
		return "", err
	}
	process.Close()

	file, err := os.CreateTemp("", "example-*.cpp")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(program)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("g++", "-std=c++17", "-O0", "-g", file.Name(), "-o", process.Name())
	_, err = cmd.Output()
	os.Remove(file.Name())

	return process.Name(), nil
}

type Var struct {
	Name string
}

type LineEvaluator struct {
	gdb    *gdb.Gdb
	output chan string
}

func New() (*LineEvaluator, error) {
	gdb, err := gdb.New(nil)
	if err != nil {
		return nil, err
	}

	output := make(chan string)
	go func() {
		for {
			line, err := readUntil(gdb, 10)
			if err == io.EOF {
				return
			}
			output <- string(line)
		}
	}()

	return &LineEvaluator{
		gdb:    gdb,
		output: output,
	}, nil
}

func mapGoToCpp(val interface{}) string {
	switch data := val.(type) {
	case int:
		return strconv.Itoa(data)
	case string:
		// TODO
		return data
	case Var:
		return data.Name
	default:
		panic("Invalid mapGoToCpp value")
	}
}

func readUntil(reader io.Reader, sep byte) ([]byte, error) {
	var res []byte
	buf := make([]byte, 1)
	for {
		_, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			return res, err
		}
		if buf[0] == sep {
			return res, nil
		}

		if buf[0] == 13 {
			continue
		}
		res = append(res, buf[0])
	}
}

func (lev *LineEvaluator) send(operation string, arguments ...string) (map[string]interface{}, error) {
	// TODO(ivo): Add Verbose option to control these logs
	log.Printf("%s %s;\n", operation, strings.Join(arguments, ", "))
	out, err := lev.gdb.Send(operation, arguments...)
	log.Printf("%v\n", out)
	return out, err
}

func (lev *LineEvaluator) Output() chan string {
	return lev.output
}

func (lev *LineEvaluator) CurrentFuncName() (string, error) {
	res, err := lev.send("stack-info-frame")
	if err != nil {
		return "", err
	}
	payload := res["payload"].(map[string]interface{})
	frame := payload["frame"].(map[string]interface{})
	return frame["func"].(string), nil
}

func (lev *LineEvaluator) CurrentLine() (int, error) {
	res, err := lev.send("stack-info-frame")
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

	_, err = lev.send("file-exec-and-symbols", program)
	if err != nil {
		return err
	}

	_, err = lev.send("break-insert", "main")
	if err != nil {
		return err
	}

	_, err = lev.send("exec-run")
	if err != nil {
		return err
	}

	return nil
}

func (lev *LineEvaluator) EvaluateFile(path string) error {
	_, err := lev.send("file-exec-and-symbols", path)
	if err != nil {
		return err
	}

	_, err = lev.send("break-insert", "main")
	if err != nil {
		return err
	}

	_, err = lev.send("exec-run")
	if err != nil {
		return err
	}

	name, _ := lev.CurrentFuncName()
	line, _ := lev.CurrentLine()
	fmt.Printf("%s %d", name, line)

	return nil
}

func (lev *LineEvaluator) EvaluateExp(exp string) (interface{}, error) {
	res, err := lev.send("data-evaluate-expression", exp)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)

	payload := res["payload"].(map[string]interface{})
	return payload["value"], nil
}

func (lev *LineEvaluator) EvaluateFunc(funcName string, signature string, values ...interface{}) (func() bool, error) {
	mappedValues := make([]string, len(values))
	for i, val := range values {
		mappedValues[i] = mapGoToCpp(val)
	}
	functionCall := fmt.Sprintf("reinterpret_cast<%s>(%s)(%s)", signature, funcName, strings.Join(mappedValues, ", "))

	_, err := lev.send("break-insert", funcName)
	if err != nil {
		return nil, err
	}

	_, err = lev.send("data-evaluate-expression", functionCall)
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

		_, err = lev.send("exec-next")
		if err != nil {
			log.Println(err)
		}
		return true
	}, nil
}

func (lev *LineEvaluator) Close() error {
	return lev.gdb.Exit()
}

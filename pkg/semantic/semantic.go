package semantic

import (
	"encoding/json"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Program struct {
	Source     string
	executable string
}

func NewProgram(source string) *Program {
	return &Program{
		Source:     source,
		executable: "",
	}
}

func (p *Program) Close() error {
	if p.executable == "" {
		return nil
	}
	return os.Remove(p.executable)
}

func (p *Program) Executable() (string, error) {
	if p.executable != "" {
		return p.executable, nil
	}

	path, err := p.compile()
	if err != nil {
		return "", err
	}
	p.executable = path
	return p.executable, nil
}

func createTemp(data string) (string, error) {
	file, err := os.CreateTemp("", "example-*.cpp")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func (p *Program) compile() (string, error) {
	file, err := os.Create("/app/runtime/program.cpp")
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.WriteString(InstrumentProgram(p.Source))

	cmd := exec.Command("cmake", "-DCMAKE_BUILD_TYPE=Debug", "-S", "/app/runtime", "-B", "/app/runtime/build/")
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	cmd = exec.Command("make", "-j4")
	cmd.Dir = "/app/runtime/build/"
	out, err := cmd.Output()
	if err != nil {
		println(out)
		return "", err
	}

	return "/app/runtime/build/runtime", nil
}

// TODO(ivo): Yea I know, it's terrible, but I've got a few days until I've got to ship it
func InstrumentProgram(program string) string {
	mods := []struct {
		*regexp.Regexp
		string
	}{
		{regexp.MustCompile(`List *\* *(\w+);`), `Pointer<List> $1("$1");`},
		{regexp.MustCompile(`(\w+) *= *new *List;`), "SetListPointerValue($1);"},
		{regexp.MustCompile(`(\w+)->val = (\w+|\d+);`), "SetListNodeValue($1, $2);"},
		{regexp.MustCompile(`(\w+)->n = (\w+|\d+);`), "SetListNodeNext($1, $2);"},
		{regexp.MustCompile(`(\w+) *= *(.+);`), "SetListPointerValue($1, $2);"},
		{regexp.MustCompile(`List *\* *(\w+)`), "Pointer<List> &$1"},
	}

	for _, mod := range mods {
		program = mod.Regexp.ReplaceAllString(program, mod.string)
	}

	lines := strings.Split(program, "\n")
	for i, _ := range lines {
		if lines[i] == "" {
			continue
		}
		if strings.Contains(lines[i], "{") {

		}
		if strings.Contains(lines[i], "#") {
			continue
		}
		if strings.Contains(lines[i], "}") {
			continue
		}
		if strings.HasPrefix(lines[i], "void") {
			continue
		}
		if strings.HasPrefix(lines[i], "using") {
			continue
		}

		//lines[i] = "AddStep();" + lines[i]
	}
	lines = append([]string{`#include "runtime.hpp"`}, lines...)
	return strings.Join(lines, "\n")
}

func (p *Program) Reflect() (*File, error) {
	temp, err := createTemp(p.Source)
	if err != nil {
		return nil, err
	}
	defer os.Remove(temp)

	cmd := exec.Command("cpp_reflect_cmd", temp)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var file File
	err = json.Unmarshal(out, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

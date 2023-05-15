package semantic

import (
	"encoding/json"
	"os"
	"os/exec"
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
	process, err := os.CreateTemp("", "program-*")
	if err != nil {
		return "", err
	}
	process.Close()

	file, err := createTemp(p.Source)
	if err != nil {
		return "", err
	}
	defer os.Remove(file)

	cmd := exec.Command("g++", "-std=c++17", "-O0", "-g", file, "-o", process.Name())
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	return process.Name(), nil
}

func (p *Program) AddRuntime() {

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

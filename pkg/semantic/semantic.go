package semantic

import (
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

func (p *Program) compile() (string, error) {
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
	defer os.Remove(file.Name())

	_, err = file.WriteString(p.Source)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("g++", "-std=c++17", "-O0", "-g", file.Name(), "-o", process.Name())
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	return process.Name(), nil
}

func (p *Program) AddRuntime() {

}

func (p *Program) Reflect() *File {

}

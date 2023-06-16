package service

import (
	"encoding/json"
	"github.com/lispberry/viz-service/pkg/evaluation"
	"github.com/lispberry/viz-service/pkg/semantic"
	"github.com/lispberry/viz-service/pkg/visualization"
)

type CompilerService struct {
	program *semantic.Program
	file    *semantic.File
	lev     *evaluation.LineEvaluator
	vis     *visualization.Visualizer
}

/* Test API !!! */

func NewCompilerService(sourceCode string) (*CompilerService, error) {
	program := semantic.NewProgram(sourceCode)

	file, err := program.Reflect()
	if err != nil {
		program.Close()
		return nil, err
	}

	lev, err := evaluation.New()
	if err != nil {
		program.Close()
		return nil, err
	}

	vis, err := visualization.NewVisualizer()
	if err != nil {
		program.Close()
		lev.Close()
		return nil, err
	}

	return &CompilerService{
		program: program,
		file:    file,
		lev:     lev,
		vis:     vis,
	}, nil
}

func (s *CompilerService) Close() {
	s.program.Close()
	s.lev.Close()
}

func (s *CompilerService) Run() error {
	path, err := s.program.Executable()
	if err != nil {
		return err
	}

	err = s.lev.EvaluateFile(path)
	if err != nil {
		return err
	}

	return nil
}

func readAll[T any](ch chan T) []T {
	var res []T
	for {
		select {
		case val, ok := <-ch:
			if ok {
				res = append(res, val)
			} else {
				return res
			}
		default:
			return res
		}
	}
}

func (s *CompilerService) RunFunction() (func() ([]string, int, bool), error) {
	fn := s.file.Functions[0]
	next, err := s.lev.EvaluateFunc(fn.Name, fn.Signature(), fn.Callable(10)...)
	if err != nil {
		return nil, err
	}

	return func() ([]string, int, bool) {
		ok := next()
		lineNumber, err := s.lev.CurrentLine()
		if err != nil {
			return nil, 0, false
		}
		if !ok {
			s.vis.DeleteStackFrame()
			return s.vis.Changes(), 0, false
		}

		lines := readAll(s.lev.Output())
		for _, line := range lines {
			var ops []visualization.RawOp
			err = json.Unmarshal([]byte(line), &ops)
			if err != nil {
				return nil, 0, false
			}
			s.vis.Apply(ops)
		}
		return s.vis.Changes(), lineNumber, true
	}, nil
}

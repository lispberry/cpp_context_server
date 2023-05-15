package main

import (
	"github.com/lispberry/viz-service/pkg/evaluation"
	"log"
)

func main() {
	lev, err := evaluation.New()
	if err != nil {
		log.Println(err)
	}
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
		log.Println(err)
	}
	_, err = lev.EvaluateFunc("program", 10, 20, 30)
	if err != nil {
		log.Println(err)
	}
	_, err = lev.EvaluateFunc("program", 10, 20, 30)
	if err != nil {
		log.Println(err)
	}
}

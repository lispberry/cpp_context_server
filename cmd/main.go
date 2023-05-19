package main

import (
	"github.com/lispberry/viz-service/pkg/visualization"
	"io"
	"log"
	"net/http"
)

func CompileHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(204)
			return
		}

		program, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(404)
			return
		}

		vis, err := visualization.NewVisualizer()
		if err != nil {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(200)
		//w.Write([]byte(res))
	}
}

func main() {
	http.HandleFunc("/api/v1/compile", CompileHandler())
	http.ListenAndServe(":8080", nil)
}

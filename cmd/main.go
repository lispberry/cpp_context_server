package main

import (
	"encoding/json"
	"fmt"
	"github.com/lispberry/viz-service/pkg/service"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Resp struct {
	Dots       []string `json:"dots"`
	LineNumber int      `json:"lineNumber"`
}

var ser *service.CompilerService
var iter func() ([]string, int, bool)

// api/v1/start
func StartHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(204)
			return
		}

		var err error
		iter, err = ser.RunFunction()
		if err != nil {
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(200)
		//w.Write([]byte(res))
	}
}

// api/v1/next
func NextHandler() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if req.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.WriteHeader(204)
			return
		}

		vis, lineNumber, hasNext := iter()
		if !hasNext {
			w.WriteHeader(200)
			return
		}
		data, err := json.Marshal(&Resp{
			Dots:       vis,
			LineNumber: lineNumber,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(200)
			return
		}
		fmt.Println(string(data))
		w.Write(data)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	ser, err = service.NewCompilerService(`void insertEnd(List * head, int val)
{
  while (head->n != nullptr)
  {
    head = head->n;
  }
  List * var;
  var = new List;
  var->val = val;
  head->n = var;
}`)
	defer ser.Close()
	if err != nil {
		return
	}
	ser.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			ser.Close()
			os.Exit(1)
		}
	}()

	http.HandleFunc("/api/v0/start", StartHandler())
	http.HandleFunc("/api/v0/next", NextHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

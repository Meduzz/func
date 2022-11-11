package main

import (
	"time"

	funclib "github.com/Meduzz/func"
	"github.com/Meduzz/wendy"
)

type (
	Greeting struct {
		Name string `json:"name"`
	}

	Reply struct {
		Greeting string `json:"greeting"`
	}
)

func main() {
	text()
	json()
}

func json() {
	// timeout of 100ms usually fails the first run of the command...
	cmd := funclib.NewCaller("./lambda", 100*time.Millisecond, "call")

	name := &Greeting{}
	name.Name = "world"

	req := wendy.Request{}
	req.Module = "test"
	req.Method = "test"
	req.Body = wendy.Json(name)
	res := &wendy.Response{}

	err := cmd.CallJSON(req, res)

	if err != nil {
		panic(err)
	}

	reply := &Reply{}
	res.Body.Bind(reply)

	println(reply.Greeting)
}

func text() {
	cmd := funclib.NewCaller("ls", 35*time.Millisecond, "-l")
	text, err := cmd.CallText("")

	if err != nil {
		panic(err)
	}

	println(text)
}

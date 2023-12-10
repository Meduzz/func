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
	// first run of a new binary takes a bit extra time
	cmd := funclib.NewCaller("./lambda", 500*time.Millisecond, "call")

	name := &Greeting{}
	name.Name = "world"

	req := wendy.Request{}
	req.Module = "test"
	req.Method = "test"
	req.Body = wendy.Json(name)
	res := &wendy.Response{}

	logz, err := cmd.CallJSON(req, res)

	println(string(logz))

	if err != nil {
		panic(err)
	}

	reply := &Reply{}
	res.Body.Bind(reply)

	println(reply.Greeting)
}

func text() {
	cmd := funclib.NewCaller("ls", 35*time.Millisecond, "-l")
	text, logz, err := cmd.CallText("")

	println(string(logz))

	if err != nil {
		panic(err)
	}

	println(text)
}

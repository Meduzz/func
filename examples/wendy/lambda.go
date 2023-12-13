package main

import (
	"fmt"
	"os"

	funclib "github.com/Meduzz/func"
	"github.com/Meduzz/helper/starters"
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
	module := wendy.NewModule("", "test")
	module.WithHandler("test", handler)

	cmds := funclib.WendyFunc(handler)
	root := starters.Root("0.1")
	root.AddCommand(cmds...)

	root.Execute()
}

func handler(r *wendy.Request) *wendy.Response {
	greeting := &Greeting{}
	r.Body.AsJson(greeting)

	fmt.Fprintln(os.Stderr, "Almost done!")

	return wendy.Ok(wendy.Json(&Reply{fmt.Sprintf("Hello %s!", greeting.Name)}))
}

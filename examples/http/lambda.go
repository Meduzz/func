package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	funclib "github.com/Meduzz/func"
	"github.com/Meduzz/helper/starters"
)

type (
	Greeting struct {
		Name string `json:"name"`
	}
)

func main() {
	lambda := funclib.HttpFunc(func(res http.ResponseWriter, req *http.Request) {
		resBody := make(map[string]interface{})

		reqBody := &Greeting{}
		decoder := json.NewDecoder(req.Body)
		decoder.Decode(reqBody)
		resBody["message"] = fmt.Sprintf("Hello %s!", reqBody.Name)

		res.Header().Add("Content-Type", "application/json")

		bs, _ := json.Marshal(resBody)
		res.Write(bs)
	})

	root := starters.Root("0.1")
	root.AddCommand(lambda)

	root.Execute()
}

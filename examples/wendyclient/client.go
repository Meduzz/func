package main

import (
	"encoding/json"
	"net"

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
	client, err := funclib.Dial("unix", "./this.sock")

	if err != nil {
		panic(err)
	}

	defer client.Close()

	var enc *json.Encoder
	var dec *json.Decoder

	client.WithConn(func(conn net.Conn) {
		dec = json.NewDecoder(conn)
		enc = json.NewEncoder(conn)
	})

	for i := 0; i < 10; i++ {
		name := &Greeting{}
		name.Name = "world"

		req := wendy.Request{}
		req.Module = "test"
		req.Method = "test"
		req.Body = wendy.Json(name)

		err = enc.Encode(req)

		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 10; i++ {
		res := &wendy.Response{}
		err = dec.Decode(res)

		if err != nil {
			panic(err)
		}

		reply := &Reply{}
		res.Body.Bind(reply)

		println(reply.Greeting)
	}

	println("Done!")
}

package funclib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Meduzz/wendy"
	"github.com/spf13/cobra"
)

// Turn a bunch of wendy modules into cobra commands.
func WendyModules(modules ...*wendy.Module) []*cobra.Command {
	cmds := make([]*cobra.Command, 0)

	handler := wendy.NewLocal(modules...)

	call := &cobra.Command{}
	call.Use = "call"
	call.Short = "reads request from stdin and writes response to stdout"
	call.RunE = wendyModuleCallHandler(handler)

	listen := &cobra.Command{}
	listen.Use = "listen"
	listen.Short = "reads request from .sock-file and writes response back to it"
	listen.Flags().String("bind", ":8080", ":port|/some.sock")
	listen.RunE = wendyModuleListenHandler(handler)

	cmds = append(cmds, call, listen)

	return cmds
}

func wendyModuleCallHandler(handler wendy.Wendy) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		bs, err := io.ReadAll(os.Stdin)

		if err != nil {
			if err != io.ErrUnexpectedEOF {
				return err
			}
		}

		req := &wendy.Request{}
		err = json.Unmarshal(bs, req)

		if err != nil {
			return err
		}

		before := time.Now()
		res := handler.Handle(req)
		after := time.Now()

		log.Printf("Func took: %s\n", after.Sub(before).String())

		bs, err = json.Marshal(res)

		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(bs)

		return err
	}
}

func wendyModuleListenHandler(handler wendy.Wendy) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		bind, err := c.Flags().GetString("bind")

		if err != nil {
			return err
		}

		network := "tcp"

		if strings.HasSuffix(bind, ".sock") {
			network = "unix"
		}

		srv, err := Listen(network, bind)

		if err != nil {
			return err
		}

		srv.Run(serveHandler(handler))

		return srv.Close()
	}
}

func serveHandler(handler wendy.Wendy) func(net.Conn) {
	return func(conn net.Conn) {
		decoder := json.NewDecoder(conn)
		encoder := json.NewEncoder(conn)

		for decoder.More() {
			req := &wendy.Request{}
			decoder.Decode(req)
			before := time.Now()
			res := handler.Handle(req)
			after := time.Now()

			log.Printf("Func took: %s\n", after.Sub(before).String())
			encoder.Encode(res)
		}

		fmt.Println("Connection closed.")
	}
}

// Turn a wendy.Handler into a bunch of cobra commands.
func WendyFunc(handler wendy.Handler) []*cobra.Command {
	cmds := make([]*cobra.Command, 0)

	call := &cobra.Command{}
	call.Use = "call"
	call.Short = "reads request from stdin and writes response to stdout"
	call.RunE = wendyMethodCallHandler(handler)

	listen := &cobra.Command{}
	listen.Use = "listen"
	listen.Short = "reads request from .sock-file and writes response back to it"
	listen.Flags().String("bind", ":8080", ":port|/some.sock")
	listen.RunE = wendyMethodListenHandler(handler)

	cmds = append(cmds, call, listen)

	return cmds
}

func wendyMethodCallHandler(method wendy.Handler) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		bs, err := io.ReadAll(os.Stdin)

		if err != nil {
			if err != io.ErrUnexpectedEOF {
				return err
			}
		}

		req := &wendy.Request{}
		err = json.Unmarshal(bs, req)

		if err != nil {
			return err
		}

		res := method(req)

		bs, err = json.Marshal(res)

		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(bs)

		return err
	}
}

func wendyMethodListenHandler(method wendy.Handler) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		bind, err := c.Flags().GetString("bind")

		if err != nil {
			return err
		}

		network := "tcp"

		if strings.HasSuffix(bind, ".sock") {
			network = "unix"
		}

		srv, err := Listen(network, bind)

		if err != nil {
			return err
		}

		go cleanup(srv)

		srv.Run(serveMethod(method))

		return nil
	}
}

func serveMethod(method wendy.Handler) func(net.Conn) {
	return func(conn net.Conn) {
		decoder := json.NewDecoder(conn)
		encoder := json.NewEncoder(conn)

		for decoder.More() {
			req := &wendy.Request{}
			decoder.Decode(req)
			before := time.Now()
			res := method(req)
			after := time.Now()

			log.Printf("Func took: %s\n", after.Sub(before).String())
			encoder.Encode(res)
		}

		fmt.Println("Connection closed.")
	}
}

func cleanup(srv *Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	srv.Close()
}

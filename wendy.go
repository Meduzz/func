package funclib

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/Meduzz/wendy"
	"github.com/spf13/cobra"
)

// Turn a bunch of wendy modules into cobra commands.
func WendyModules(app string, modules ...*wendy.Module) []*cobra.Command {
	cmds := make([]*cobra.Command, 0)

	handler := wendy.NewLocal(app, modules...)

	call := &cobra.Command{}
	call.Use = "call"
	call.Short = "reads request from stdin and writes response to stdout"
	call.RunE = wendyModuleCallHandler(handler)

	cmds = append(cmds, call)

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

		ctx := context.Background()
		res := handler.Handle(ctx, req)

		bs, err = json.Marshal(res)

		if err != nil {
			return err
		}

		_, err = os.Stdout.Write(bs)

		return err
	}
}

// Turn a wendy.Handler into a bunch of cobra commands.
func WendyFunc(handler wendy.Handler) []*cobra.Command {
	cmds := make([]*cobra.Command, 0)

	call := &cobra.Command{}
	call.Use = "call"
	call.Short = "reads request from stdin and writes response to stdout"
	call.RunE = wendyMethodCallHandler(handler)

	cmds = append(cmds, call)

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

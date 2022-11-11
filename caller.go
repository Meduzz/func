package funclib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type (
	Caller struct {
		cmd     string
		args    []string
		env     []string
		timeout time.Duration
	}
)

// Create a new caller
func NewCaller(cmd string, timeout time.Duration, args ...string) *Caller {
	return &Caller{
		cmd:     cmd,
		args:    args,
		timeout: timeout,
	}
}

// Add something to the env of the command.
func (c *Caller) AddEnv(key, value string) {
	c.env = append(c.env, fmt.Sprintf("%s=%s", key, value))
}

// Call the command and optionally send json to stdin and expect json from stdout.
func (c *Caller) CallJSON(in interface{}, out interface{}) error {
	var bs []byte
	var err error
	if in != nil {
		bs, err = json.Marshal(in)

		if err != nil {
			return err
		}
	}

	cmd, cancel := c.command()

	defer cancel()

	writer, _ := cmd.StdinPipe()
	stdoutSync := c.stdout(cmd)

	err = cmd.Start()

	if err != nil {
		return err
	}

	if bs != nil {
		_, err = writer.Write(bs)

		if err != nil {
			cmd.Process.Kill()
			return err
		}
	}

	writer.Close()

	buf, err := stdoutSync()

	if err != nil {
		cmd.Process.Kill()
		return err
	}

	err = cmd.Wait()

	if err != nil {
		return err
	}

	if cmd.ProcessState.Success() {
		err = json.Unmarshal(buf.Bytes(), out)

		if err != nil {
			return err
		}
	}

	return nil
}

// Call the command and optionally send text to stdin and read text from stdout.
func (c *Caller) CallText(in string) (string, error) {
	cmd, cancel := c.command()

	defer cancel()

	writer, _ := cmd.StdinPipe()
	stdoutSync := c.stdout(cmd)

	err := cmd.Start()

	if err != nil {
		return "", err
	}

	if in != "" {
		_, err = writer.Write([]byte(in))

		if err != nil {
			cmd.Process.Kill()
			return "", err
		}
	}

	writer.Close()

	buf, err := stdoutSync()

	if err != nil {
		cmd.Process.Kill()
		return "", err
	}

	err = cmd.Wait()

	if err != nil {
		return "", err
	}

	if cmd.ProcessState.Success() {
		return buf.String(), nil
	}

	// TODO weird spot!
	return "", nil
}

func (c *Caller) command() (*exec.Cmd, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)

	return exec.CommandContext(ctx, c.cmd, c.args...), cancel
}

func (c *Caller) stdout(cmd *exec.Cmd) func() (*bytes.Buffer, error) {
	reader, _ := cmd.StdoutPipe()
	buffer := bytes.NewBufferString("")

	return func() (*bytes.Buffer, error) {
		_, err := buffer.ReadFrom(reader)

		return buffer, err
	}
}

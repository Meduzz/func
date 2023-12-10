package funclib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
func (c *Caller) CallJSON(in interface{}, out interface{}) ([]byte, error) {
	var bs []byte
	var err error
	if in != nil {
		bs, err = json.Marshal(in)

		if err != nil {
			return []byte(err.Error()), err
		}
	}

	cmd, cancel := c.command()

	defer cancel()

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	stdoutSync := c.outputCollector(stdout)
	stderrSync := c.outputCollector(stderr)

	err = cmd.Start()

	if err != nil {
		return []byte(err.Error()), err
	}

	if bs != nil {
		_, err = stdin.Write(bs)

		if err != nil {
			cmd.Process.Kill()
			return []byte(err.Error()), err
		}
	}

	stdin.Close()

	buf, err := stdoutSync()

	if err != nil {
		cmd.Process.Kill()
		return []byte(err.Error()), err
	}

	errorz, err := stderrSync()

	if err != nil {
		cmd.Process.Kill()
		return []byte(err.Error()), err
	}

	err = cmd.Wait()

	if err != nil {
		return errorz.Bytes(), err
	}

	if cmd.ProcessState.Success() {
		if out != nil {
			err = json.Unmarshal(buf.Bytes(), out)

			if err != nil {
				return errorz.Bytes(), err
			}
		}
	}

	return errorz.Bytes(), nil
}

// Call the command and optionally send text to stdin and read text from stdout.
func (c *Caller) CallText(in string) (string, []byte, error) {
	cmd, cancel := c.command()

	defer cancel()

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	stdoutSync := c.outputCollector(stdout)
	stderrSync := c.outputCollector(stderr)

	err := cmd.Start()

	if err != nil {
		return "", nil, err
	}

	if in != "" {
		_, err = stdin.Write([]byte(in))

		if err != nil {
			cmd.Process.Kill()
			return "", nil, err
		}
	}

	stdin.Close()

	buf, err := stdoutSync()

	if err != nil {
		cmd.Process.Kill()
		return "", nil, err
	}

	errorz, err := stderrSync()

	if err != nil {
		cmd.Process.Kill()
		return "", nil, err
	}

	err = cmd.Wait()

	if err != nil {
		return "", errorz.Bytes(), err
	}

	if cmd.ProcessState.Success() {
		return buf.String(), errorz.Bytes(), nil
	}

	// TODO weird spot!
	return "", errorz.Bytes(), nil
}

func (c *Caller) command() (*exec.Cmd, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)

	return exec.CommandContext(ctx, c.cmd, c.args...), cancel
}

func (c *Caller) outputCollector(reader io.ReadCloser) func() (*bytes.Buffer, error) {
	buffer := bytes.NewBufferString("")

	return func() (*bytes.Buffer, error) {
		_, err := buffer.ReadFrom(reader)

		return buffer, err
	}
}

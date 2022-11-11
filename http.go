package funclib

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/spf13/cobra"
)

// HttpFunc turns a HandleFunc into a callable lambda.
func HttpFunc(handler http.HandlerFunc) *cobra.Command {
	call := &cobra.Command{}
	call.Use = "call"
	call.Short = "reads the request from stdin and writes the response to stdout"
	call.RunE = httpFunctionCallWrapper(handler)

	return call
}

func httpFunctionCallWrapper(handler http.HandlerFunc) func(*cobra.Command, []string) error {
	return func(c *cobra.Command, s []string) error {
		req := bytes.NewBufferString("")
		_, err := io.Copy(req, os.Stdin)

		if err != nil {
			if err != io.ErrUnexpectedEOF {
				return err
			}
		}

		httpReq, err := http.ReadRequest(bufio.NewReader(req))

		if err != nil {
			return err
		}

		httpRes := httptest.NewRecorder()

		handler(httpRes, httpReq)

		return httpRes.Result().Write(os.Stdout)
	}
}

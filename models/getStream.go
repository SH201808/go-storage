package models

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type GetStream struct {
	reader io.Reader
}

func newGetStream(url string, object string) (*GetStream, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "NewGetStream NewRequest: hash:%s, url: %s", object, url)
	}
	req.Header.Set("object", object)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", resp.StatusCode)
	}
	return &GetStream{resp.Body}, nil
}

func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	url := "http://" + server + "/file/download"
	return newGetStream(url, object)
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

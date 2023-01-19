package models

import (
	"encoding/json"
	response "file-server/models/Response"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type TempPutStraem struct {
	UUID   string
	Server string
}

func NewPutStream(Size, Hash, Server string) (*TempPutStraem, error) {
	url := "http://" + Server + "/temp/fileMeta"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "NewPutStream NewRequest: hash:%s, size: %s, ip: %s", Hash, Server, Size)
	}
	req.Header.Set("Hash", Hash)
	req.Header.Set("Size", Size)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "NewPutStream Do req: hash:%s, size: %s, ip: %s", Hash, Server, Size)
	}
	data, err := io.ReadAll(resp.Body)
	respData := &response.Response{}
	json.Unmarshal(data, &respData)
	if err != nil {
		return nil, errors.Wrapf(err, "NewPutStream Read Body: hash:%s, size: %s, ip: %s", Hash, Server, Size)
	}
	uuid := respData.Params["uuid"]
	return &TempPutStraem{
		Server: Server,
		UUID:   uuid.(string),
	}, nil
}

func (stream *TempPutStraem) Write(data []byte) (n int, err error) {
	url := "http://" + stream.Server + "/temp/file/"
	req, err := http.NewRequest("PATCH", url, strings.NewReader(string(data)))
	if err != nil {
		return 0, errors.Wrapf(err, "upload temp file err NewRequest: uuid:", stream.UUID)
	}
	req.Header.Set("uuid", stream.UUID)

	client := http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return 0, errors.Wrapf(err, "upload temp file err client Do: uuid:", stream.UUID)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, errors.Errorf("Write:dataServer return not ok:%d", resp.StatusCode)
	}
	return len(data), nil
}

func (stream *TempPutStraem) Commit(flag bool) {
	method := "PUT"
	url := "http://" + stream.Server + "/temp"
	if !flag {
		method = "DELETE"
		url += "/fileDelete"
	}
	url += "/removeToStore"
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("uuid", stream.UUID)

	client := http.Client{}
	client.Do(req)
}

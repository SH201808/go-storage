package rs

import (
	"encoding/base64"
	"encoding/json"
	"file-server/models"
	"file-server/utils"
	"fmt"
	"io"
	"log"
	"net/http"
)

type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	UUIDS   []string
}

type ResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

func NewResumablePutStream(dataServers []string, name, hash string, size int64) (*ResumablePutStream, error) {
	putStream, err := NewRsPutStream(dataServers, hash, size)
	if err != nil {
		return nil, err
	}
	uuids := make([]string, ALL_SHARDS)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*models.TempPutStraem).UUID
	}
	token := &resumableToken{Name: name, Hash: hash, Size: size, Servers: dataServers, UUIDS: uuids}
	return &ResumablePutStream{putStream, token}, nil
}

func NewRSResumablePutStreamFromToken(token string) (*ResumablePutStream, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var t resumableToken
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &models.TempPutStraem{Server: t.Servers[i], UUID: t.UUIDS[i]}
	}
	enc := NewEncoder(writers)
	return &ResumablePutStream{&RSPutStream{enc}, &t}, nil
}

func (s *ResumablePutStream) CurrentSize() int64 {
	r, err := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.UUIDS[0]))
	if err != nil {
		log.Println(err)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := utils.GetOffsetFromHeader(r.Header) * DATA_SHARDS
	if size > s.Size {
		size = s.Size
	}
	return size
}

func (s *ResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	// todo 对token加密
	return base64.StdEncoding.EncodeToString(b)
}

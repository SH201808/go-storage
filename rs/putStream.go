package rs

import (
	"file-server/models"
	"fmt"
	"io"
	"log"

	"github.com/klauspost/reedsolomon"
)

type RSPutStream struct {
	*encoder
}

func NewRsPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	writers := make([]io.Writer, ALL_SHARDS)
	var err error
	for i := range writers {
		// todo rs 的最后一块不足 BLOCK_PER_SHARD 会补0
		if i == ALL_SHARDS-1 && perShard%BLOCK_PER_SHARD != 0 {

		}
		writers[i], err = models.NewPutStream(fmt.Sprintf("%d", perShard), fmt.Sprintf("%s.%d", hash, i), dataServers[i])
		if err != nil {
			return nil, err
		}
	}
	enc := NewEncoder(writers)

	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(success bool) error {
	s.Flush()
	for i := range s.writers {
		err := s.writers[i].(*models.TempPutStraem).Commit(success)
		if err != nil {
			return err
		}
	}
	return nil
}

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder
	cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &encoder{writers, enc, nil}
}

func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	for current := 0; length != 0; {
		next := BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		if len(e.cache) == BLOCK_SIZE {
			err := e.Flush()
			if err != nil {
				log.Println("flush err:", err)
				return length, nil
			}
		}
		current += next
		length -= next
	}
	return len(p), nil
}

func (e *encoder) Flush() error {
	if len(e.cache) == 0 {
		return nil
	}
	shards, _ := e.enc.Split(e.cache)
	err := e.enc.Encode(shards)
	if err != nil {
		log.Println("encode err:", err)
		return err
	}
	for i := range shards {
		n, err := e.writers[i].Write(shards[i])
		if err != nil {
			log.Println("encoder write n:", n)
			log.Println("encoder write err:", err)
			return err
		}
	}
	e.cache = []byte{}
	return nil
}

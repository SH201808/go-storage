package models

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value float32
	}
}

package models

import "strings"

type TempFileMeta struct {
	UUID string
	Name string
	Size string
}

func (t *TempFileMeta) Hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *TempFileMeta) Id() string {
	s := strings.Split(t.Name, ".")
	return s[1]
}

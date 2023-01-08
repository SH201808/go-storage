package UUID

import uuid "github.com/satori/go.uuid"

func Gen() string {
	return uuid.NewV4().String()
}

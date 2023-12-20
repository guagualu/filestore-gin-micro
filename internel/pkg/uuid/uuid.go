package uuid

import (
	guuid "github.com/google/uuid"
)

func NewUuid() string {
	return guuid.New().String()
}

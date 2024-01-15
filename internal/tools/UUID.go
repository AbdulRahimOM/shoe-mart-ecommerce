package tools

import "github.com/google/uuid"

func MakeRandomUUID() (string, error) {
	uuidValue, err := uuid.NewRandom()
	return uuidValue.String(), err
}

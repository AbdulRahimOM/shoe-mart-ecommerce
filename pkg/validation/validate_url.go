package requestuestValidation

import (
	"strconv"
	"strings"
)

func ValidateAndParseID(id string) (int, error) {
	if _, err := strconv.Atoi(id); err != nil {
		return 0, err
	}
	return strconv.Atoi(id)
}
func ValidateAndParseIDs(param string) ([]uint, error) {
	var ids []uint

	idStrings := strings.Split(param, ",")

	for _, idStr := range idStrings {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, uint(id))
	}

	return ids, nil
}

func ValidateAndParseInt(param string) (int, error) {
	var paramInt int
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, err

	}
	return paramInt, nil
}

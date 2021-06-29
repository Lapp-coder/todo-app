package request

import (
	"errors"
	"strings"
)

const (
	defaultTime = "19700101:"
)

func parseCompletedDate(date *string) error {
	strDate := *date

	index := strings.Index(strDate, ":")
	if index == -1 {
		return errors.New("failed to parse completed date")
	}

	*date = strDate[:index] + " " + strDate[index+1:]

	return nil
}

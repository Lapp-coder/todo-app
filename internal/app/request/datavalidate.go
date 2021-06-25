package request

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func getTimeNow() string {
	year := strconv.Itoa(time.Now().Year())
	mouth := time.Now().Month().String()
	day := strconv.Itoa(time.Now().Day())

	return year + mouth + day
}

func parseCompletedDate(date *string) error {
	strDate := *date

	index := strings.Index(strDate, ":")
	if index == -1 {
		return errors.New("failed to parse completed date")
	}

	*date = strDate[:index] + " " + strDate[index+1:]

	return nil
}

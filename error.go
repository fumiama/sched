package sched

import (
	"strconv"
	"strings"
)

// Errors parallel errors with batch index.
type Errors []error

func (errs Errors) Error() string {
	sb := strings.Builder{}
	for i, err := range []error(errs) {
		if err == nil {
			continue
		}
		sb.WriteByte('#')
		sb.WriteString(strconv.Itoa(i))
		sb.WriteByte(':')
		sb.WriteString(err.Error())
		sb.WriteByte(' ')
	}
	return sb.String()
}

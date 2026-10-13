package sqlite

import (
	"fmt"
	"strings"
)

func inQuery[T any](format string, ids[] T) (string, []any) {
	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	return fmt.Sprintf(format, strings.Join(placeholders, ",")), args
}

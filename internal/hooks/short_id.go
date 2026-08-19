package hooks

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

func AssignShortID(e *core.RecordEvent) error {
	if e.Record.GetInt("short_id") == 0 {
		var maximum int
		if err := e.App.RecordQuery("sounds").
			Select("COALESCE(MAX(short_id), 0)").
			Row(&maximum); err != nil {
			return fmt.Errorf("failed to pick next short_id: %w", err)
		}
		e.Record.Set("short_id", maximum+1)
	}

	return e.Next()
}

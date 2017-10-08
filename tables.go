package mn_pg

import (
	"database/sql"
	"fmt"
	"strings"
)

func ManageTables(db *sql.DB, toDrop []string, toCreate [][2]string) error {
	if len(toDrop) > 0 {
		dq := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", strings.Join(toDrop, ", "))
		if _, err := db.Exec(dq); err != nil {
			return fmt.Errorf("table drop failed: %v", err)
		}
	}
	if len(toCreate) == 0 {
		return nil
	}
	for _, d := range toCreate {
		if _, err := db.Exec(d[1]); err != nil {
			return fmt.Errorf("table %s create failed: %v", d[0], err)
		}
	}
	return nil
}

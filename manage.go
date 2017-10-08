package mn_pg

import (
	"fmt"
	"io"
	"strings"
)

type Opts struct {
	DBName string `json:"dbName"`
	DBUser string `json:"dbUser"`
	DBPWD  string `json:"dbPassword"`

	SchemaFileNames []string `json:"schemaFilenames"`

	DropAll  bool     `json:"defaultDropAll"`
	DropSome []string `json:"defaultDropSome"`

	CreateAll  bool     `json:"defaultCreateAll"`
	CreateSome []string `json:"defaultCreateSome"`
}

func Manage(w io.Writer, j *Opts, schema [][2]string) (err error) {
	if schema == nil {
		schema, err = LoadSchemaFromFiles(j.SchemaFileNames...)
		if err != nil {
			return err
		}
	}
	var toDrop []string
	var toCreate [][2]string
	sMap := make(map[string][2]string, len(schema))
	for _, s := range schema {
		if j.DropAll {
			toDrop = append(toDrop, s[0])
		}
		if j.CreateAll {
			toCreate = append(toCreate, s)
		} else {
			sMap[s[0]] = s
		}
	}
	if !j.DropAll {
		toDrop = j.DropSome
	}
	if !j.CreateAll {
		for _, t := range j.CreateSome {
			if s, ok := sMap[t]; ok {
				toCreate = append(toCreate, s)
			} else {
				return fmt.Errorf("no schema found for table %s", t)
			}
		}
	}
	if len(toDrop) == 0 && len(toCreate) == 0 {
		return fmt.Errorf("to tables found to drop/create")
	}

	db, err := LoadDB(j.DBUser, j.DBPWD, j.DBName)
	if err != nil {
		return err
	}
	defer db.Close()
	err = ManageTables(db, toDrop, toCreate)
	if w != nil && err == nil {
		if len(toDrop) > 0 {
			fmt.Fprintf(w, "Dropped tables: %s\n", strings.Join(toDrop, ", "))
		}
		if len(toCreate) > 0 {
			parts := make([]string, len(toCreate))
			for i, s := range toCreate {
				parts[i] = s[0]
			}
			fmt.Fprintf(w, "Created tables: %s\n", strings.Join(parts, ", "))
		}
	}
	return err
}

package mn_pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
)

func LoadDB(user, pass, dbName string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, pass, dbName))
	if err != nil {
		return nil, fmt.Errorf("loaddb failure: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("loaddb ping failure: %v", err)
	}
	return db, nil
}

func LoadSchemaFromFiles(fileNames ...string) ([][2]string, error) {
	usedNames := map[string]string{}
	var schema [][2]string
	for _, fileName := range fileNames {
		sql, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("schema file %s error: %v", fileName, err)
		}
		parsed, err := ParseSQL(sql)
		if err != nil {
			return nil, fmt.Errorf("schema file %s error: %v", fileName, err)
		}
		for _, s := range parsed {
			if f2, ok := usedNames[s[0]]; ok {
				return nil, fmt.Errorf("multiple definition of table %s in files %s and %s", s[0], f2, fileName)
			} else {
				usedNames[s[0]] = fileName
				schema = append(schema, s)
			}
		}

	}
	return schema, nil
}

func LoadOptsFromJSON(fileName string) (*Opts, error) {
	d := new(Opts)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return d, json.NewDecoder(file).Decode(d)
}

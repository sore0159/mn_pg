package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func ManageTables(w io.Writer, db *sql.DB, opts *Options) error {
	tables, err := ParseSchema(opts.SchemaFileNames...)
	if err != nil {
		return err
	}
	if opts.DropFlag {
		var drop []string
		if len(opts.ToDrop) == 0 {
			for _, t := range tables {
				drop = append(drop, t[0])
			}
		} else {
			drop = opts.ToDrop
		}
		dq := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", strings.Join(drop, ", "))
		if _, err := db.Exec(dq); err != nil {
			return fmt.Errorf("table drop failed: %v", err)
		} else if w != nil {
			fmt.Fprintf(w, "Dropped tables %s\n", strings.Join(drop, ", "))
		}
	}
	if !opts.CreateFlag {
		return nil
	}
	var create [][2]string
	if len(opts.ToCreate) == 0 {
		create = tables
	} else {
		for _, t := range opts.ToCreate {
			var found bool
			for _, test := range tables {
				if test[0] == t {
					found = true
					create = append(create, test)
					break
				}
			}
			if !found {
				return fmt.Errorf("no sql found for table %s", t)
			}
		}
	}
	for _, d := range create {
		if _, err := db.Exec(d[1]); err != nil {
			return fmt.Errorf("table %s create failed: %v", d[0], err)
		} else if w != nil {
			fmt.Fprintf(w, "Created table %s\n", d[0])
		}
	}
	return nil
}

func ParseSchema(fileNames ...string) ([][2]string, error) {
	tables := [][2]string{}
	usedNames := map[string]string{}
	for _, fileName := range fileNames {
		sql, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("schema file %s error: %v", fileName, err)
		}
		var open bool
		for i, ln := range bytes.Split(sql, []byte("\n")) {
			ln = bytes.TrimSpace(ln)
			if len(ln) == 0 {
				continue
			}
			if bytes.HasPrefix(bytes.ToLower(ln), []byte("create table ")) {
				if open {
					return nil, fmt.Errorf("create table during open definition on line %d", i)
				}

				flds := bytes.Fields(ln)
				if len(flds) < 3 {
					return nil, fmt.Errorf("can't find table name on line %d", i)
				}
				name := string(flds[2])
				if fName2, ok := usedNames[name]; ok {
					if fName2 == fileName {
						return nil, fmt.Errorf("multiple definition of table %s in file %s", name, fileName)
					} else {
						return nil, fmt.Errorf("multiple definition of table %s in files %s and %s", name, fName2, fileName)
					}
				}
				usedNames[name] = fileName
				tbl := [2]string{
					name,
					string(ln),
				}
				tables = append(tables, tbl)
				open = true
			} else if !open {
				return nil, fmt.Errorf("create table write to unkown table line %d", i)
			} else {
				if bytes.HasPrefix(ln, []byte(");")) {
					open = false
				}
				tbl := &tables[len(tables)-1]
				tbl[1] += string(ln)
			}
		}
		if open {
			return nil, fmt.Errorf("unclosed sql schema file")
		}
	}
	return tables, nil
}

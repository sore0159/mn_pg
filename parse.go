package mn_pg

import (
	"bytes"
	"fmt"
)

func ParseSQL(sql []byte) ([][2]string, error) {
	tables := [][2]string{}
	usedNames := map[string]bool{}

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
			if usedNames[name] {
				return nil, fmt.Errorf("multiple definition of table %s", name)
			}
			usedNames[name] = true
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
	return tables, nil
}

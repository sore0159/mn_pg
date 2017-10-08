package main

import (
	"log"
	"os"
)

func main() {
	opts, err := ParseOptions()
	if err != nil {
		log.Printf("Error in command parsing: %v\n", err)
		if opts != nil && opts.HelpFlag {
			PrintHelp()
		}
		return
	}
	if opts.HelpFlag {
		PrintHelp()
		return
	}
	db, err := LoadDB(opts.DBUser, opts.DBPWD, opts.DBName)
	if err != nil {
		log.Printf("Error loading DB %s: %v\n", opts.DBName, err)
		return
	}
	defer db.Close()
	if err := ManageTables(os.Stdout, db, opts); err != nil {
		log.Printf("Error managing DB %s: %v\n", opts.DBName, err)
		return
	}
}

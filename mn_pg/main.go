package main

import (
	"log"
	"os"

	mnpg "github.com/sore0159/mn_pg"
)

func main() {
	opts, err := ParseCommandLine()
	if err != nil {
		log.Printf("Error in commandline parsing: %v\n", err)
		if opts != nil && opts.HelpFlag {
			PrintHelp()
		}
		return
	}
	if opts.HelpFlag {
		PrintHelp()
		return
	}
	j, err := LoadJSON(opts)
	if err != nil {
		log.Printf("Error loading json configuration: %v\n", err)
		PrintHelp()
		return
	} else if opts.HelpFlag {
		PrintHelp()
		return
	}
	err = mnpg.Manage(os.Stdout, j, nil)
	if err != nil {
		log.Printf("Error during management: %v\n", err)
		PrintHelp()
		return
	}
}

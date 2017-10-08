package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HelpFlag bool

	JSONFile string

	DropFlag   bool
	ToDrop     []string
	CreateFlag bool
	ToCreate   []string
	TestFlag   bool
}

func ParseCommandLine() (*Config, error) {
	o := new(Config)
	var state int
	const DFLAG = 1
	const CFLAG = 2
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h":
			o.HelpFlag = true
			return o, nil
		case "-d":
			state = DFLAG
			o.DropFlag = true
		case "-c":
			state = CFLAG
			o.CreateFlag = true
		default:
			switch state {
			case DFLAG:
				o.ToDrop = append(o.ToDrop, arg)
			case CFLAG:
				o.ToCreate = append(o.ToCreate, arg)
			default:
				if strings.HasPrefix(arg, "-") {
					for _, char := range arg {
						switch char {
						case '-':
						case 'h':
							o.HelpFlag = true
						case 'd':
							o.DropFlag = true
						case 'c':
							o.CreateFlag = true
						default:
							o.HelpFlag = true
							return o, fmt.Errorf("unknown flag '%s'", string(char))
						}
					}
				} else if strings.HasSuffix(arg, ".json") {
					if o.JSONFile == "" {
						o.JSONFile = arg
					} else {
						o.HelpFlag = true
						return o, fmt.Errorf("multiple options json files not supported")
					}
				} else {
					o.HelpFlag = true
					return o, fmt.Errorf("unknown arg %s", arg)
				}
			}
		}
	}
	if o.JSONFile == "" {
		o.HelpFlag = true
		return o, nil
	}
	return o, nil
}

func PrintHelp() {
	fmt.Printf(`
mn_pg JSON_FILENAME [-h][-d [DROP_NAMES...]][-c [CREATE_NAMES...]]
  -h	Print this help message and do nothing else
  -d	Drop tables listed, all present in schema files if none listed
  -c	Create tables listed, all present in schema files if none listed
You may combine flags that have no arguments (e.g. 'manage -cd')

The first argument must end in '.json' and mn_pg will attempt to load database configurations from that file.  All command line arguments will will supercede options in this file.

`)
}

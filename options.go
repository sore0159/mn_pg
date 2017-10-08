package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type FileOpts struct {
	DBName string `json:"dbName"`
	DBUser string `json:"dbUser"`
	DBPWD  string `json:"dbPassword"`

	SchemaFileNames []string `json:"schemaFilenames"`

	DropAll  bool     `json:"defaultDropAll"`
	DropSome []string `json:"defaultDropSome"`

	CreateAll  bool     `json:"defaultCreateAll"`
	CreateSome []string `json:"defaultCreateSome"`
}

func LoadFileOpts(fileName string) (*FileOpts, error) {
	d := new(FileOpts)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return d, json.NewDecoder(file).Decode(d)
}

type Options struct {
	HelpFlag bool

	DBName string
	DBUser string
	DBPWD  string

	SchemaFileNames []string

	DropFlag   bool
	ToDrop     []string
	CreateFlag bool
	ToCreate   []string
	TestFlag   bool
}

func (o *Options) ApplyDefaults(d *FileOpts) {
	o.DBName = d.DBName
	o.DBUser = d.DBUser
	o.DBPWD = d.DBPWD

	o.SchemaFileNames = d.SchemaFileNames

	if !o.DropFlag && (d.DropAll || len(d.DropSome) > 0) {
		o.DropFlag = true
		if !d.DropAll {
			o.ToDrop = d.DropSome
		}
	}
	if !o.CreateFlag && (d.CreateAll || len(d.CreateSome) > 0) {
		o.CreateFlag = true
		if !d.CreateAll {
			o.ToCreate = d.CreateSome
		}
	}
}

func ParseOptions() (*Options, error) {
	o := new(Options)
	var d *FileOpts
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
					if d == nil {
						var err error
						d, err = LoadFileOpts(arg)
						if err != nil {
							o.HelpFlag = true
							if os.IsNotExist(err) {
								return o, fmt.Errorf("cannot find options file %s")
							} else {
								return o, fmt.Errorf("cannot parse options file %s: %v", arg, err)
							}
						}
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
	if d == nil {
		o.HelpFlag = true
		return o, nil
	}
	o.ApplyDefaults(d)
	if (!o.DropFlag && !o.CreateFlag) || len(o.SchemaFileNames) == 0 {
		o.HelpFlag = true
		return o, nil
	}
	return o, nil
}

func PrintHelp() {
	fmt.Printf(`
mn_pg OPTIONS_FILENAME [-h][-d [DROP_NAMES...]][-c [CREATE_NAMES...]]
  -h	Print this help message and do nothing else
  -d	Drop tables listed, all present in schema files if none listed
  -c	Create tables listed, all present in schema files if none listed
You may combine flags that have no arguments (e.g. 'manage -cd')

  The first argument must end in '.json' and mn_pg will attempt to load database related options from that file.  All command line specifications will will superceed options in this file.

`)
}

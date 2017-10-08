package main

import (
	mnpg "github.com/sore0159/mn_pg"
)

func LoadJSON(cfg *Config) (*mnpg.Opts, error) {
	j, err := mnpg.LoadOptsFromJSON(cfg.JSONFile)
	if err != nil {
		return nil, err
	}

	if cfg.DropFlag {
		if len(cfg.ToDrop) == 0 {
			j.DropAll = true
		} else {
			j.DropAll = false
			j.DropSome = cfg.ToDrop
		}
	}
	if cfg.CreateFlag {
		if len(cfg.ToCreate) == 0 {
			j.CreateAll = true
		} else {
			j.CreateAll = false
			j.CreateSome = cfg.ToCreate
		}
	}
	return j, nil
}

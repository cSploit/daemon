package main

import "flag"
import "github.com/ianschenck/envflag"

var configPath string

func init() {
	const (
		configPathDefault = "./config.json"
		configPathUsage   = "path to config file"
	)

	env_prefix := "CSPLOIT_"

	flag.StringVar(&configPath, "config", configPathDefault, configPathUsage)
	flag.StringVar(&configPath, "c", configPathDefault, configPathUsage)
	envflag.StringVar(&configPath, env_prefix+"CONFIG", configPathDefault, configPathUsage)
}

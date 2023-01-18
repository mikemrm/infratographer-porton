package main

import (
	"github.com/infratographer/porton/plugin"
)

const (
	// PluginName is the name of the plugin.
	// This is the name that will be used to identify the plugin in the configuration.
	PluginName = "porton"
)

var (
	// Implement symbols that the plugin loader will look for.
	HandlerRegisterer = plugin.NewPortonRegisterer(PluginName)
	ClientRegisterer  = plugin.NewPortonRegisterer(PluginName)
)

func main() {}

package main

import (
	"github.com/infratographer/porton/plugin"
)

var (
	// Implement symbols that the plugin loader will look for.
	ModifierRegisterer = plugin.NewPortonRegisterer(plugin.PluginName)
)

func main() {}

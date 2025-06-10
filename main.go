package main

import (
	scaffolder2 "github.com/SwanHtetAungPhyo/gostart/scaffolder"
	"github.com/SwanHtetAungPhyo/gostart/wizzard"
	"os"

	"github.com/fatih/color"
)

func main() {
	wizard := wizzard.NewWizard()
	config, err := wizard.Run()
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	scaffolder := scaffolder2.NewScaffolder(config)
	if err := scaffolder.CreateProject(); err != nil {
		color.Red("Error creating project: %v", err)
		os.Exit(1)
	}

	color.Green("\n✅ Project '%s' created successfully!", config.ProjectDir)
	color.Cyan("📁 Next steps:")
	color.Yellow("   cd %s\n", config.ProjectDir)
	if config.UseAir {
		color.Red("Please try to fix the air.toml setting to get the hot reload. Because, air init generate the default setting")
		color.Yellow("In this version, please run with the makefile command , \n make build")
		color.Blue("   air ")
	} else {
		color.Green("   go run ./cmd")
	}
}

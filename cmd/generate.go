package cmd

import (
	"fmt"
	"log"

	"github.com/devops-kung-fu/common/util"
	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/devops-kung-fu/gardener/lib"
)

var (
	deflate     bool
	generateCmd = &cobra.Command{
		Use:     "generate [path]",
		Example: "gardener generate .",
		Aliases: []string{"gen"},
		Short:   "Generates links in any Markdown files that are pointed to PlantUML source files",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				fmt.Println()
				log.Fatal("no path was passed to gardener")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			generate(args[0])
		},
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().BoolVar(&deflate, "deflate", true, "compresses the generated URLs")
}

func generate(path string) {
	diagramFiles, e := lib.FindFiles(Afs, path, ".*\\.(pu|puml|plantuml|iuml|wsd)")
	if e != nil {
		log.Fatal(e)
	}
	util.DoIf(Verbose, func() {
		color.Style{color.FgLightBlue, color.OpBold}.Print("Generating Links...\n\n")
		util.PrintInfo(fmt.Sprintf("Found %x diagrams", len(diagramFiles)))
		for _, file := range diagramFiles {
			util.PrintTabbed(file)
		}
		util.PrintInfo("Processing Markdown files")
	})

	_, _ = lib.Generate(Afs, diagramFiles, path, Verbose, deflate)

	util.DoIf(Verbose, func() {
		util.PrintSuccess("Done!\n")
	})

}

package cmd

import (
	"fmt"
	"log"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/devops-kung-fu/gardener/src"
	"github.com/devops-kung-fu/gardener/util"
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
	diagramFiles, e := src.FindFiles(Afs, path, ".*\\.(pu|puml|plantuml|iuml|wsd)")
	if e != nil {
		log.Fatal(e)
	}
	util.DoIf(Verbose, func() {
		color.Style{color.FgLightBlue, color.OpBold}.Print("Generating Links...\n\n")
		util.PrintInfo(fmt.Sprintf("Found %x diagrams", len(diagramFiles)))
		util.PrintInfo("Processing Markdown files")
	})

	markdownFiles, e := src.FindFiles(Afs, path, ".*\\.md")
	if e != nil {
		log.Fatal(e)
	}
	for _, markdownFile := range markdownFiles {
		util.DoIf(Verbose, func() {
			log.Print("Working on ", markdownFile)
			util.PrintTabbed(markdownFile)
		})

		for _, diagramFile := range diagramFiles {
			diagramContent := src.ReadFileContentBytes(Afs, diagramFile)
			var url string
			if deflate {
				log.Print("Deflate Encoding Diagram for: ", diagramFile)
				url = src.DeflateEncodedURL(diagramContent)
			} else {
				log.Print("Hex Encoding Diagram for: ", diagramFile)
				url = src.HexEncodedURL(diagramContent)
			}

			searchImageStub := fmt.Sprintf("\\!\\[%s\\]\\(.*\\)", diagramFile)
			replaceImageStub := fmt.Sprintf("![%s](%s)", diagramFile, url)
			src.ReplaceLineInFile(Afs, markdownFile, searchImageStub, replaceImageStub)
		}
	}
	util.DoIf(Verbose, func() {
		util.PrintSuccess("Done!\n")
	})

}

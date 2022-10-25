// Package cmd contains all of the commands that may be executed in the cli
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/devops-kung-fu/common/util"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	Afs     = &afero.Afero{Fs: afero.NewOsFs()}
	version = "1.0.1"
	debug   bool
	Verbose bool
	rootCmd = &cobra.Command{
		Use:     "gardener",
		Short:   `Converts PlantUML source files into links in Markdown`,
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !debug {
				log.SetOutput(ioutil.Discard)
			}
			util.DoIf(Verbose, func() {
				fmt.Println()
				color.Style{color.FgWhite, color.OpBold}.Println("█▀▀ ▄▀█ █▀█ █▀▄ █▀▀ █▄ █ █▀▀ █▀█")
				color.Style{color.FgWhite, color.OpBold}.Println("█▄█ █▀█ █▀▄ █▄▀ ██▄ █ ▀█ ██▄ █▀▄")
				fmt.Println()
				fmt.Println("https://github.com/devops-kung-fu/gardener")
				fmt.Printf("Version: %s\n", version)
				fmt.Println()
			})
		},
	}
)

// Execute creates the command tree and handles any error condition returned
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "show debug output")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", true, "show verbose output")
}

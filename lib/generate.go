package lib

import (
	"fmt"
	"log"

	"github.com/devops-kung-fu/common/util"
	"github.com/spf13/afero"
)

func Generate(afs *afero.Afero, diagramFiles []string, path string, verbose bool, deflate bool) (markdownFiles []string, err error) {
	markdownFiles, err = FindFiles(afs, path, ".*\\.md")
	if err != nil {
		log.Fatal(err)
	}
	for _, markdownFile := range markdownFiles {
		util.DoIf(verbose, func() {
			log.Print("Working on ", markdownFile)
			util.PrintTabbed(markdownFile)
		})

		for _, diagramFile := range diagramFiles {
			diagramContent, err := ReadFileContentBytes(afs, diagramFile)
			if util.IsErrorBool(err) {
				log.Fatal(err)
			}
			var url string
			if deflate {

				log.Print("Deflate Encoding Diagram for: ", diagramFile)
				url = DeflateEncodedURL(diagramContent)
			} else {
				log.Print("Hex Encoding Diagram for: ", diagramFile)
				url = HexEncodedURL(diagramContent)
			}

			searchImageStub := fmt.Sprintf("\\!\\[%s\\]\\(.*\\)", diagramFile)
			replaceImageStub := fmt.Sprintf("![%s](%s)", diagramFile, url)
			_, _ = ReplaceLineInFile(afs, markdownFile, searchImageStub, replaceImageStub)
		}
	}
	return
}

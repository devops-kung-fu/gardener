package src

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/afero"
)

const ENCODING = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
const PLANTUML_URL = "http://www.plantuml.com/plantuml/png"

func plantUMLBase64(input []byte) []byte {
	encoding := base64.NewEncoding(ENCODING)
	return []byte(encoding.EncodeToString(input))
}

// func brotliCompress(content []byte) ([]byte, error) {
// 	var comp bytes.Buffer
// 	w := brotli.NewWriterLevel(&comp, 11)
// 	_, writeErr := w.Write(content)
// 	if writeErr != nil {
// 		return nil, writeErr
// 	}
// 	_ = w.Close()
// 	return comp.Bytes(), nil
// }

func deflateCompress(content []byte) ([]byte, error) {
	var comp bytes.Buffer
	w, _ := flate.NewWriter(&comp, flate.HuffmanOnly)
	_, writeErr := w.Write(content)
	if writeErr != nil {
		return nil, writeErr
	}
	_ = w.Flush()
	_ = w.Close()
	return comp.Bytes(), nil
}

func FindFiles(afs *afero.Afero, path string, re string) ([]string, error) {
	libRegEx, e := regexp.Compile(re)
	if e != nil {
		return nil, e
	}

	var files []string
	e = afs.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			files = append(files, filePath)
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return files, nil
}

func ReadFileContentString(afs *afero.Afero, path string) string {
	content, err := afs.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func ReadFileContentBytes(afs *afero.Afero, filePath string) []byte {
	content, err := afs.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func HexEncodedURL(content []byte) string {
	encodedStr := hex.EncodeToString(content)
	return fmt.Sprintf("%s/~h%s", PLANTUML_URL, encodedStr)
}

func DeflateEncodedURL(content []byte) string {
	comp, err := deflateCompress(content)
	if err != nil {
		log.Fatal()
	}
	encoded := plantUMLBase64(comp)
	return fmt.Sprintf("%s/%s", PLANTUML_URL, string(encoded))
}

// func brotliEncodedURL(content []byte) string {
// 	comp, err := brotliCompress(content)
// 	if err != nil {
// 		log.Fatal()
// 	}
// 	encoded := plantUMLBase64(comp)
// 	return fmt.Sprintf("%s/0%s", PLANTUML_URL, string(encoded))
// }

func ReplaceLineInFile(afs *afero.Afero, filePath string, searchString string, replaceString string) bool {
	libRegEx, e := regexp.Compile(searchString)
	if e != nil {
		log.Fatal(e)
	}
	content := ReadFileContentString(afs, filePath)
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if libRegEx.MatchString(line) {
			libRegEx.ReplaceAllString(line, replaceString)
			lines[i] = replaceString
		}
	}
	output := strings.Join(lines, "\n")
	err := afs.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.Print("Updated image tags in ", filePath)
	return true
}

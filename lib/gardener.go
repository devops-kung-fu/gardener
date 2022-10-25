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

func deflateCompress(content []byte) ([]byte, error) {
	var comp bytes.Buffer
	w, _ := flate.NewWriter(&comp, flate.HuffmanOnly)
	_, _ = w.Write(content)
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
		} else {
			return err
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return files, e
}

func ReadFileContentString(afs *afero.Afero, path string) (string, error) {
	content, err := afs.ReadFile(path)
	if err != nil {
		log.Println(err)
		return string(""), err
	}
	return string(content), err
}

func ReadFileContentBytes(afs *afero.Afero, filePath string) ([]byte, error) {
	content, err := afs.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return content, err
}

func HexEncodedURL(content []byte) string {
	encodedStr := hex.EncodeToString(content)
	return fmt.Sprintf("%s/~h%s", PLANTUML_URL, encodedStr)
}

func DeflateEncodedURL(content []byte) string {
	comp, err := deflateCompress(content)
	if err != nil {
		log.Println()
	}
	encoded := plantUMLBase64(comp)
	return fmt.Sprintf("%s/%s", PLANTUML_URL, string(encoded))
}

func ReplaceLineInFile(afs *afero.Afero, filePath string, searchString string, replaceString string) (bool, error) {
	libRegEx, e := regexp.Compile(searchString)
	if e != nil {
		log.Println(e)
		return false, e
	}
	content, err := ReadFileContentString(afs, filePath)
	if err != nil {
		return false, err
	}
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if libRegEx.MatchString(line) {
			libRegEx.ReplaceAllString(line, replaceString)
			lines[i] = replaceString
		}
	}
	output := strings.Join(lines, "\n")
	_ = afs.WriteFile(filePath, []byte(output), 0644)
	log.Print("Updated image tags in ", filePath)
	return true, nil
}

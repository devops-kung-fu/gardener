package src

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFindFiles(t *testing.T) {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}
	files, err := FindFiles(afs, ".", ".*\\.(pu|puml|plantuml|iuml|wsd)")
	assert.NoError(t, err)
	assert.Len(t, files, 0)

	afs.WriteFile("test.pu", []byte("test"), 0644)
	files, err = FindFiles(afs, ".", ".*\\.(pu|puml|plantuml|iuml|wsd)")

	assert.NoError(t, err)
	assert.Len(t, files, 1)
}

func TestReadFileContentString(t *testing.T) {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}
	afs.WriteFile("example.pu", []byte("test"), 0644)
	result := ReadFileContentString(afs, "example.pu")

	assert.Equal(t, "test", result)
}

func TestReadFileContentBytes(t *testing.T) {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}
	afs.WriteFile("example.pu", []byte("test"), 0644)
	result := ReadFileContentBytes(afs, "example.pu")

	assert.Equal(t, []byte("test"), result)
}

func TestHexEncodedURL(t *testing.T) {
	result := HexEncodedURL([]byte(`Hello World!`))
	assert.Equal(t, "http://www.plantuml.com/plantuml/png/~h48656c6c6f20576f726c6421", result)
}

func TestReplaceLineInFile(t *testing.T) {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}
	afs.WriteFile("README.md", []byte("\\!\\[example.pu\\]\\(.*\\)"), 0644)
	result := ReplaceLineInFile(afs, "README.md", "\\!\\[example.pu\\]\\(.*\\)", "![example.pu](https://example.com)")

	assert.True(t, result)
}

func TestDeflateEncodedURL(t *testing.T) {
	result := DeflateEncodedURL([]byte("@startuml\nAlice -> Bob: Authentication Request\nBob --> Alice: Authentication Response\nAlice -> Bob: Another authentication Request\nAlice <-- Bob: Another authentication Response\n@enduml"))
	assert.Equal(t, "http://www.plantuml.com/plantuml/png/1C3HZSCW40JGVwgO1cZ0Ebd19RW3u4PY9RARcA7_lDTIVRJVCvLfdSWdhcW7ojQWotgLXUFcTtCfNT6GyuaohVD0sHfqMQ-oSDnSd_35LAPr8f-ueXqe7XfyKBS6NTQhB1mtlvjBgKphn5_EkA8TA1uQV52t1btMgomSDzSdV36zwF_xFNy0003__m400F__", result)
}

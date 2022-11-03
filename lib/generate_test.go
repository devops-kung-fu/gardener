package lib

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}
	afs.WriteFile("example.pu", samplePlantUMLFile(), 0644)
	afs.WriteFile("README.md", []byte("![example.pu]()"), 0644)

	diagramFiles, err := FindFiles(afs, ".", ".*\\.(pu|puml|plantuml|iuml|wsd)")

	assert.NoError(t, err)
	result, err := Generate(afs, diagramFiles, ".", true, true)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	file, _ := afs.ReadFile("README.md")
	assert.Equal(t, file, []byte("![example.pu](http://www.plantuml.com/plantuml/png/1C3XQSGm30NWt_8KRm4l4Cfn3dI1T0B7VM218-LY6Jf-loJodb6VMDV0Zrxt8Bx_wdKF9f4oj17vXTtF3ML5fuMs6kg6Wv_56KbaznCvlr26DcueApejjDLGDnoSjzjaIY9bQ2Fo2xkV6ufvT3weApejjDLGDnpyBCv88cLe8xvonkxdXc8UdG_gYauBRJNK3GVdxJQPHDwqttMD9Fy0003__m400F__)"))
}

func samplePlantUMLFile() []byte {
	test := `
		@startuml Simple Example
		Alice -> Bob: Authentication Request
		Bob --> Alice: Authentication Response

		Alice -> Bob: Another authentication Request
		Alice <-- Bob: Another authentication Response
		@enduml
	`
	return []byte(test)
}

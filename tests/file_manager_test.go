// file: src/data/file_manager_test.go
package data

import (
	"os"
	"testing"

	"github.com/utkarsh5026/Genka/src/data"
)

func TestSaveLangFiles(t *testing.T) {
	fm, err := data.NewFileManager()
	if err != nil {
		t.Fatalf("Failed to create FileManager: %v", err)
	}

	langs := []data.Language{data.LangEnglish, data.LangSpanish}
	data := [][]byte{
		[]byte(`{"hello": "world"}`),
		[]byte(`{"hola": "mundo"}`),
	}

	filePaths, err := fm.SaveLangFiles(langs, data)
	if err != nil {
		t.Fatalf("Failed to save language files: %v", err)
	}

	for lang, path := range filePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("File for language %s was not created at path %s", lang, path)
		}
	}
}

func TestSaveDataFiles(t *testing.T) {
	fm, err := data.NewFileManager()
	if err != nil {
		t.Fatalf("Failed to create FileManager: %v", err)
	}

	fileNames := []data.FileName{data.ArtifactMainStatFile, data.ArchonDataFile}
	data := [][]byte{
		[]byte(`{"score": 100}`),
		[]byte(`{"score": 200}`),
	}

	filePaths, err := fm.SaveDataFiles(fileNames, data)
	if err != nil {
		t.Fatalf("Failed to save data files: %v", err)
	}

	for fileName, path := range filePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("File for data %s was not created at path %s", fileName, path)
		}
	}
}

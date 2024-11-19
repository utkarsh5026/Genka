package data

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// createDefaultDirectoryPath returns the default directory path for data files
func createDefaultDirectoryPath() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to get current directory")
	}
	dirPath := filepath.Join(filepath.Dir(filename), "..", "..", "data")

	if err := validatePath(dirPath, true); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	return dirPath, nil
}

type FileManager struct {
	directoryPath string
	langPath      string
	dataPath      string
}

func NewFileManager() (*FileManager, error) {
	dirPath, err := createDefaultDirectoryPath()
	if err != nil {
		return nil, err
	}

	fmt.Println("Saving files to", dirPath)
	return &FileManager{
		directoryPath: dirPath,
		langPath:      filepath.Join(dirPath, "langs"),
		dataPath:      filepath.Join(dirPath, "data"),
	}, nil
}

func (fm *FileManager) SaveLangFiles(langs []Language, data [][]byte) (map[Language]string, error) {
	langFilePaths := make(map[Language]string)
	if err := validatePath(fm.langPath, true); err != nil {
		return nil, fmt.Errorf("failed to create lang directory: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var saveErr error

	for i, lang := range langs {
		wg.Add(1)
		go func(lang Language) {
			defer wg.Done()
			filePath := filepath.Join(fm.langPath, fmt.Sprintf("%s.json", lang))
			if err := os.WriteFile(filePath, data[i], 0644); err != nil {
				mu.Lock()
				saveErr = fmt.Errorf("failed to save lang file: %w", err)
				mu.Unlock()
				return
			}
			mu.Lock()
			langFilePaths[lang] = filePath
			mu.Unlock()
		}(lang)
	}
	wg.Wait()

	if saveErr != nil {
		return nil, saveErr
	}
	return langFilePaths, nil
}

// validatePath checks if a directory path exists and optionally creates it if it doesn't.
// If createIfNotExists is true and the path doesn't exist, it will create the directory
// with 0755 permissions. Returns an error if the path check fails or if directory
// creation fails when requested.
func validatePath(path string, createIfNotExists bool) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if createIfNotExists {
			return os.MkdirAll(path, 0755)
		}
		return err
	}
	return nil
}

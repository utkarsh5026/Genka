package data

import (
	"encoding/json"
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

// saveFiles saves multiple files to the specified path with the given names and data.
// It creates the directory if it doesn't exist and saves files concurrently.
//
// Parameters:
//   - path: The directory path where files will be saved
//   - names: Slice of file names (without .json extension)
//   - data: Slice of byte slices containing the file contents
//
// Returns:
//   - map[string]string: Map of file names to their full file paths
//   - error: Any error that occurred during saving
func (fm *FileManager) saveFiles(path string, names []string, data [][]byte) (map[string]string, error) {
	filePaths := make(map[string]string)
	if err := validatePath(path, true); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var saveErr error

	for i, name := range names {
		wg.Add(1)
		go func(i int, name string) {
			defer wg.Done()
			filePath := filepath.Join(path, fmt.Sprintf("%s.json", name))

			if !isValidJSON(data[i]) {
				mu.Lock()
				saveErr = fmt.Errorf("invalid JSON data for file: %s", name)
				mu.Unlock()
				return
			}

			if err := os.WriteFile(filePath, data[i], 0644); err != nil {
				mu.Lock()
				saveErr = fmt.Errorf("failed to save file: %w", err)
				mu.Unlock()
				return
			}
			mu.Lock()
			filePaths[name] = filePath
			mu.Unlock()
		}(i, name)
	}
	wg.Wait()

	if saveErr != nil {
		return nil, saveErr
	}
	return filePaths, nil
}

// SaveLangFiles saves language files to the language directory.
// It takes a slice of Language enums and their corresponding data,
// saves them as JSON files, and returns a map of Languages to file paths.
//
// Parameters:
//   - langs: Slice of Language enums representing the languages to save
//   - data: Slice of byte slices containing the language file contents
//
// Returns:
//   - map[Language]string: Map of Languages to their saved file paths
//   - error: Any error that occurred during saving
func (fm *FileManager) SaveLangFiles(langs []Language, data [][]byte) (map[Language]string, error) {
	langNames := make([]string, len(langs))
	for i, lang := range langs {
		langNames[i] = string(lang)
	}
	langFilePaths, err := fm.saveFiles(fm.langPath, langNames, data)
	if err != nil {
		return nil, err
	}

	// Convert map[string]string to map[Language]string
	result := make(map[Language]string)
	for i, path := range langFilePaths {
		result[Language(i)] = path
	}
	return result, nil
}

// SaveDataFiles saves game data files to the data directory.
// It takes a slice of FileName enums and their corresponding data,
// saves them as JSON files, and returns a map of FileNames to file paths.
//
// Parameters:
//   - filesNames: Slice of FileName enums representing the files to save
//   - data: Slice of byte slices containing the file contents
//
// Returns:
//   - map[FileName]string: Map of FileNames to their saved file paths
//   - error: Any error that occurred during saving
func (fm *FileManager) SaveDataFiles(filesNames []GenshinDataFileName, data [][]byte) (map[GenshinDataFileName]string, error) {
	fileNames := make([]string, len(filesNames))
	for i, fileName := range filesNames {
		fileNames[i] = string(fileName)
	}
	dataFilePaths, err := fm.saveFiles(fm.dataPath, fileNames, data)
	if err != nil {
		return nil, err
	}

	// Convert map[string]string to map[FileName]string
	result := make(map[GenshinDataFileName]string)
	for i, path := range dataFilePaths {
		result[GenshinDataFileName(i)] = path
	}
	return result, nil
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

// IsValidJSON checks if a byte slice contains valid JSON data.
// It attempts to unmarshal the data into a temporary interface{} to verify syntax.
//
// Parameters:
//   - data: Byte slice containing the potential JSON data to validate
//
// Returns:
//   - bool: True if the data is valid JSON, false otherwise
func isValidJSON(data []byte) bool {
	var js interface{}
	return json.Unmarshal(data, &js) == nil
}

func (fm *FileManager) LoadFile(file GenshinDataFileName) ([]byte, error) {
	filePath := filepath.Join(fm.dataPath, fmt.Sprintf("%s.json", file))
	return os.ReadFile(filePath)
}

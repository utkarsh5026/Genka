package data

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Language string

const (
	GitlabUrl           = "https://gitlab.com/Dimbreath/AnimeGameData/-/tree/master/"
	LanguageMapFilesUrl = "https://gitlab.com/Dimbreath/AnimeGameData/-/raw/master/TextMap/"
	GenshinDataFilesUrl = "https://gitlab.com/Dimbreath/AnimeGameData/-/raw/master/ExcelBinOutput/"

	LangSimplifiedChinese  Language = "chs"
	LangTraditionalChinese Language = "cht"
	LangGerman             Language = "de"
	LangEnglish            Language = "en"
	LangSpanish            Language = "es"
	LangFrench             Language = "fr"
	LangIndonesian         Language = "id"
	LangJapanese           Language = "jp"
	LangKorean             Language = "kr"
	LangPortuguese         Language = "pt"
	LangRussian            Language = "ru"
	LangThai               Language = "th"
	LangVietnamese         Language = "vi"
)

type ResourceLoader struct {
	fm             *FileManager
	loggingEnabled bool
	logger         *log.Logger
}

func NewResourceLoader(fm *FileManager, loggingEnabled bool) *ResourceLoader {
	var logger *log.Logger
	if loggingEnabled {
		logger = log.New(os.Stdout, "ResourceLoader: ", log.LstdFlags)
	}
	return &ResourceLoader{
		fm:             fm,
		loggingEnabled: loggingEnabled,
		logger:         logger,
	}
}

// LoadLangFiles concurrently downloads language files for all configured languages.
// It uses goroutines to fetch files in parallel, collects any errors that occur,
// and saves the downloaded files using the FileManager.
//
// The function creates an HTTP client and launches a goroutine for each language
// to download its corresponding file. It waits for all downloads to complete
// before checking for errors and saving the files.
//
// Returns:
//   - error: Returns nil if all files were successfully downloaded and saved,
//     or an error describing what went wrong during the process
func (rl *ResourceLoader) LoadLangFiles(langs []Language) error {
	var wg sync.WaitGroup
	result := make([][]byte, len(langs))
	errs := make([]error, len(langs))
	client := &http.Client{}

	// Launch goroutines for each language
	for i := range langs {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			lang := langs[idx]
			url := fmt.Sprintf("%sTextMap%s.json?ref_type=heads&inline=false",
				LanguageMapFilesUrl,
				strings.ToUpper(string(lang)),
			)
			data, err := rl.loadFileFromUrl(url, client)
			result[idx] = data
			errs[idx] = err
		}(i)
	}

	wg.Wait()

	// Check for any errors
	for i, err := range errs {
		if err != nil {
			return fmt.Errorf("failed to load lang file %s: %w", langs[i], err)
		}
	}

	_, err := rl.fm.SaveLangFiles(langs, result)
	if err != nil {
		return fmt.Errorf("failed to save lang files: %w", err)
	}

	return nil
}

// LoadDataFiles concurrently downloads game data files from the configured repository.
// It uses goroutines to fetch files in parallel, collects any errors that occur,
// and saves the downloaded files using the FileManager.
//
// Parameters:
//   - dataFiles: A slice of GenshinDataFileName values specifying which files to download
//
// The function creates an HTTP client and launches a goroutine for each data file
// to download its corresponding JSON file. It waits for all downloads to complete
// before checking for errors and saving the files.
//
// Returns:
//   - error: Returns nil if all files were successfully downloaded and saved,
//     or an error describing what went wrong during the process
func (rl *ResourceLoader) LoadDataFiles(dataFiles []GenshinDataFileName) error {
	var wg sync.WaitGroup
	result := make([][]byte, len(dataFiles))
	errs := make([]error, len(dataFiles))
	client := &http.Client{}

	for i := range dataFiles {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			file := dataFiles[idx]
			url := getDataFileUrl(file)
			data, err := rl.loadFileFromUrl(url, client)
			result[idx] = data
			errs[idx] = err
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			return fmt.Errorf("failed to load data file %s: %w", dataFiles[i], err)
		}
	}

	_, err := rl.fm.SaveDataFiles(dataFiles, result)
	if err != nil {
		return fmt.Errorf("failed to save data files: %w", err)
	}

	return nil
}

// loadFileFromUrl downloads and returns the contents of a file from the given URL.
//
// Parameters:
//   - url: The URL to download the file from
//   - client: The HTTP client to use for the request
//
// Returns:
//   - []byte: The contents of the downloaded file
//   - error: nil if successful, otherwise an error describing what went wrong
func (rl *ResourceLoader) loadFileFromUrl(url string, client *http.Client) ([]byte, error) {
	var result []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error downloading data file: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(resp.Body)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading lang file: %w", err)
	}
	result = append(result, data...)
	return result, nil
}

// GetFile loads a data file from disk or downloads it if missing.
// It first attempts to load the file from the local filesystem using FileManager.
// If the file doesn't exist and downloadIfMissing is true, it will download
// the file from the remote URL and save it locally before returning the contents.
//
// Parameters:
//   - file: The FileName enum indicating which data file to load
//   - downloadIfMissing: Whether to download the file if it doesn't exist locally
//
// Returns:
//   - []byte: The contents of the loaded file
func (rl *ResourceLoader) GetFile(file GenshinDataFileName, downloadIfMissing bool) ([]byte, error) {
	data, err := rl.fm.LoadFile(file)
	if err != nil && os.IsNotExist(err) && downloadIfMissing {
		url := getDataFileUrl(file)
		fmt.Printf("File is missing so downloading the data file %s\n from %s\n", file, url)

		data, err = rl.loadFileFromUrl(url, http.DefaultClient)
		if err != nil {
			return nil, err
		}
		_, err = rl.fm.SaveDataFiles([]GenshinDataFileName{file}, [][]byte{data})
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load data file %s: %w", file, err)
	}

	fmt.Printf("Loaded data file %s\n", file)
	return data, nil
}

// GetLangDirPath returns the path to the language files directory.
func (rl *ResourceLoader) GetLangDirPath() string {
	return rl.fm.langPath
}

// GetDataDirPath returns the path to the data files directory.
func (rl *ResourceLoader) GetDataDirPath() string {
	return rl.fm.dataPath
}

// DownloadAllDataFiles concurrently downloads all Genshin Impact data files from the remote repository and saves them locally
//
// The function spawns a goroutine for each file to download, allowing parallel downloads.
//
// Returns:
//   - error: The first error encountered during downloads, or nil if all downloads succeed
func (rl *ResourceLoader) DownloadAllDataFiles() error {
	fileNames := GetGenshinDataFileNames()
	return rl.LoadDataFiles(fileNames)
}

func (rl *ResourceLoader) DownLoadAllLanguageFiles() error {
	langs := []Language{
		LangSimplifiedChinese,
		LangTraditionalChinese,
		LangGerman,
		LangEnglish,
		LangSpanish,
		LangFrench,
		LangIndonesian,
		LangJapanese,
		LangKorean,
		LangPortuguese,
		LangRussian,
		LangThai,
		LangVietnamese,
	}

	return rl.LoadLangFiles(langs)
}

// getDataFileUrl constructs the URL for downloading a Genshin Impact data file
func getDataFileUrl(file GenshinDataFileName) string {
	return fmt.Sprintf("%s%s.json?ref_type=heads&inline=false", GenshinDataFilesUrl, file)
}

package data

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type Language string

const (
	GitlabUrl           = "https://gitlab.com/Dimbreath/AnimeGameData/-/tree/master/"
	LanguageMapFilesUrl = "https://gitlab.com/Dimbreath/AnimeGameData/-/raw/master/TextMap/"
	DataFilesUrl        = "https://gitlab.com/Dimbreath/AnimeGameData/-/raw/master/ExcelBinOutput/"

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

type ResourceLoaderOptions struct {
	DataFiles []FileName
	Langs     []Language
}

type ResourceLoader struct {
	fm   *FileManager
	opts ResourceLoaderOptions
}

func NewResourceLoader(fm *FileManager, opts ResourceLoaderOptions) *ResourceLoader {
	return &ResourceLoader{
		fm:   fm,
		opts: opts,
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
func (rl *ResourceLoader) LoadLangFiles() error {
	var wg sync.WaitGroup
	langs := rl.opts.Langs
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

// LoadDataFiles concurrently downloads game data files configured in ResourceLoaderOptions.
// It uses goroutines to fetch files in parallel, collects any errors that occur,
// and saves the downloaded files using the FileManager.
//
// The function creates an HTTP client and launches a goroutine for each data file
// to download its corresponding JSON file. It waits for all downloads to complete
// before checking for errors and saving the files.
//
// Returns:
//   - error: Returns nil if all files were successfully downloaded and saved,
//     or an error describing what went wrong during the process
func (rl *ResourceLoader) LoadDataFiles() error {
	var wg sync.WaitGroup
	dataFiles := rl.opts.DataFiles
	result := make([][]byte, len(dataFiles))
	errs := make([]error, len(dataFiles))
	client := &http.Client{}

	for i := range dataFiles {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			file := dataFiles[idx]
			url := fmt.Sprintf("%s%s.json?ref_type=heads&inline=false", DataFilesUrl, file)
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

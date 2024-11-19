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
	fm    *FileManager
	langs []Language
}

func NewResourceLoader(fm *FileManager, langs []Language) *ResourceLoader {
	return &ResourceLoader{
		fm:    fm,
		langs: langs,
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
	result := make([][]byte, len(rl.langs))
	errs := make([]error, len(rl.langs))
	client := &http.Client{}

	// Launch goroutines for each language
	for i := range rl.langs {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			data, err := rl.loadLangFile(rl.langs[idx], client)
			result[idx] = data
			errs[idx] = err
		}(i)
	}

	wg.Wait()

	// Check for any errors
	for i, err := range errs {
		if err != nil {
			return fmt.Errorf("failed to load lang file %s: %w", rl.langs[i], err)
		}
	}

	_, err := rl.fm.SaveLangFiles(rl.langs, result)
	if err != nil {
		return fmt.Errorf("failed to save lang files: %w", err)
	}

	return nil
}

// loadLangFile downloads and reads a language file for the specified language.
// It constructs the URL using the language code, makes an HTTP GET request,
// and returns the file contents as a byte slice.
//
// Parameters:
//   - lang: The Language identifier for the file to load
//   - client: The HTTP client to use for the request
//
// Returns:
//   - []byte: The contents of the language file
//   - error: Any error that occurred during the download/read process
func (rl *ResourceLoader) loadLangFile(lang Language, client *http.Client) ([]byte, error) {
	var result []byte
	url := fmt.Sprintf("%sTextMap%s.json?ref_type=heads&inline=false",
		LanguageMapFilesUrl,
		strings.ToUpper(string(lang)))

	fmt.Println("Loading lang file from", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error downloading lang file: %w", err)
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

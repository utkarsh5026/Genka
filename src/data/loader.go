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

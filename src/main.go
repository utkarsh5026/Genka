package main

import (
	"fmt"
	"github.com/utkarsh5026/Genka/src/data"
)

func main() {
	fmt.Println("Hello, World!")
	fm, err := data.NewFileManager()

	if err != nil {
		fmt.Println(err)
	}
	loader := data.NewResourceLoader(fm, []data.Language{data.LangEnglish})
	err = loader.LoadLangFiles()

	if err != nil {
		fmt.Println(err)
	}
}

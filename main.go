package main

import (
	"flag"

	"github.com/Ghytro/sme/helpers"
)

func main() {
	smeFilesDir := flag.String("smeFilesDir", "", "Directory with smep files")
	outLang := flag.String("outLang", "", "Language to generate the code")
	outDir := flag.String("outDir", "", "Where to generate the out code")
	flag.Parse()

	helpers.HandleSmeFilesDirArgumentErrors(smeFilesDir)
	helpers.HandleOutLangArgumentErrors(outLang)
	helpers.HandlerOutDirArgumentErrors(outDir)
}

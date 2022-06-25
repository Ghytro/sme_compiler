package parser

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/Ghytro/sme/helpers"
)

func Parse(smeFilesDir *string, outLang *string) {
	if err := filepath.Walk(
		*smeFilesDir,
		smeFileHandler,
	); err != nil {
		helpers.PrintError(err.Error())
	}
}

func smeFileHandler(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	ParseReaderContent(bufio.NewReader(file))
	return nil
}

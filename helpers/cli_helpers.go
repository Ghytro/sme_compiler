package helpers

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const DefaultSmeDir = "./sme"

var GeneratedLanguages = [...]string{"cpp", "java", "go", "python"}

func PathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func PrintWarning(warning string) (int, error) {
	if warning == "" {
		return 0, nil
	}
	return fmt.Fprintln(os.Stdout, "Warning: "+warning+"\n")
}

func PrintError(err string) {
	if err == "" {
		return
	}
	fmt.Fprint(os.Stderr, "Error: "+err+"\n")
	os.Exit(1)
}

func IsGeneratedLanguage(lang string) bool {
	for _, l := range GeneratedLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

func HandleSmeFilesDirArgumentErrors(smeFilesDir *string) {
	if *smeFilesDir == "" {
		PrintWarning(
			fmt.Sprintf(
				"no path specified for .sme files, using default %s directory instead",
				DefaultSmeDir,
			),
		)
		*smeFilesDir = DefaultSmeDir
	}
	exists, err := PathExists(*smeFilesDir)
	if err != nil {
		log.Printf("Debug: error in helpers.HandleSmeFilesDirArgumentsErrors %s\n", err)
		PrintError("incorrect path for dir with .sme files specified")
	}
	if !exists {
		PrintError("incorrect path specified for directory with sme files")
	}
}

func HandleOutLangArgumentErrors(outLang *string) {
	if *outLang == "" {
		PrintError(
			fmt.Sprintf(
				"language not specified, the allowed values are: %s",
				strings.Join(GeneratedLanguages[:], ", "),
			),
		)
	}
	if !IsGeneratedLanguage(*outLang) {
		PrintError(
			fmt.Sprintf(
				"incorrect value of generated language, the allowed values are: %s",
				strings.Join(GeneratedLanguages[:], ", "),
			),
		)
	}
}

func HandlerOutDirArgumentErrors(outDir *string) {
	if *outDir == "" {
		PrintError("no path specified for outgoing compiled files")
	}
	exists, err := PathExists(*outDir)
	if err != nil {
		log.Printf("Debug: error in helpers.HandleOutDirArgumentsErrors %s\n", err)
		PrintError("incorrect path for dir with outgoing compiled files")
	}
	if !exists {
		if err := os.Mkdir(*outDir, os.ModePerm); err != nil {
			log.Printf("Debug: error in helpers.HandleOutDirArgumentsErrors %s\n", err)
			PrintError(fmt.Sprintf("unable to create directory %s", *outDir))
		}
	}
}

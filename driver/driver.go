package driver

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"martinpetr.dev/kina/compiler/internal/diagnostics"
	"martinpetr.dev/kina/compiler/internal/lexer"
	"martinpetr.dev/kina/compiler/internal/performance"
	"martinpetr.dev/kina/compiler/internal/treebuilder"
	"martinpetr.dev/kina/compiler/projectConfig"
)

type Options struct {
	Out string
}

type parseProjectFilesResult struct{}

type parseFileResult struct {
	Imports []string
}

func Compile(projectPath string, opt Options) error {
	// Resolve and parse kina.toml project config file
	configPath := path.Join(projectPath, "kina.toml")
	config, err := projectConfig.ParseFile(configPath)
	if err != nil {
		return err
	}

	fmt.Printf("Compiling project %s...\n", projectPath)

	// Get the absolute path to the entry file
	absEntrypointPath, err := filepath.Abs(path.Join(projectPath, config.Project.Entry))
	if err != nil {
		return err
	}

	// Check if the entry file exists
	if _, err := os.Stat(absEntrypointPath); os.IsNotExist(err) {
		return fmt.Errorf("Entry file '%s' does not exist", absEntrypointPath)
	}

	diagnosticsBag := diagnostics.NewBag(true)

	// Parse (lex + ast) all project files
	// Resolved imports are also parsed recursively (results cached)
	_, err = parseProjectFiles(projectPath, absEntrypointPath, diagnosticsBag)
	if err != nil {
		return err
	}

	// Check if there are any diagnostics (errors/warnings) and print them
	err = diagnosticsBag.Err(os.Stderr)
	if err != nil {
		return err
	}

	return nil
}

// Parses (lex + ast) all project files starting from the entrypoint file
// Resolved imports are also parsed recursively (results cached)
func parseProjectFiles(projectRootPath string, absEntrypointPath string, diagnosticsBag *diagnostics.Bag) (*parseProjectFilesResult, error) {
	var pathsToParse []string = []string{absEntrypointPath}                             // List of files to parse (entrypoint + imports)
	var parsedFileResults map[string]parseFileResult = make(map[string]parseFileResult) // Cache of parsed file results (key is abs file path)

	var cwd, err = os.Getwd()
	if err != nil {
		return nil, err
	}

	// Loop through all files to parse
	for len(pathsToParse) > 0 {
		relativeFilePath, err := filepath.Rel(cwd, pathsToParse[0])
		if err != nil {
			return nil, err
		}

		// Read the file content
		src, err := os.ReadFile(pathsToParse[0])
		if err != nil {
			return nil, err
		}

		// Create a new diagnostics file for the current file and save it to the diagnostics files map
		diagnosticFile := diagnostics.NewFile(
			relativeFilePath,
			src,
		)

		// Parse the first file in the list
		res, err := parseFile(pathsToParse[0], src, diagnosticsBag.For(diagnosticFile))
		if err != nil {
			return nil, err
		}

		// Cache the result of the parsed file
		parsedFileResults[pathsToParse[0]] = *res

		// Add all imports of the parsed file to the list of files to parse
		// except for the ones that have already been parsed
		for _, imp := range res.Imports {
			if _, ok := parsedFileResults[imp]; !ok {
				pathsToParse = append(pathsToParse, imp)
			}
		}

		// Remove the parsed file from the list of files to parse
		pathsToParse = pathsToParse[1:]
	}

	return &parseProjectFilesResult{}, nil
}

// Parses (lex + ast) a single file and returns the result
func parseFile(filePath string, src []byte, reporter *diagnostics.Reporter) (*parseFileResult, error) {
	fmt.Printf("Parsing file %s...\n", filePath)

	startLexer := time.Now()
	performance.ReportHeapSize("start")

	// Lex the file
	lexerResult := lexer.ProcessFile(filePath, src, reporter)
	lexerTime := time.Since(startLexer)
	performance.ReportHeapSize("lexer")

	startAsi := time.Now()

	asiResult := lexerResult.InsertSemicolons()
	asiTime := time.Since(startAsi)
	performance.ReportHeapSize("asi")

	essentialTokens := asiResult.RemoveNonEssential()

	startAst := time.Now()
	tree := treebuilder.BuildTree(filePath, essentialTokens.Tokens, reporter)
	astTime := time.Since(startAst)
	performance.ReportHeapSize("ast")

	_, err := tree.String()
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(json))
	fmt.Printf("Lexer: %s, ASI: %s, AST: %s\n", lexerTime, asiTime, astTime)

	return &parseFileResult{
		Imports: []string{},
	}, nil
}

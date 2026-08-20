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
	EmitTokens bool
	EmitTree bool
}

type parseProjectFilesResult struct{}

type parseFileResult struct {
	Imports []string
	Tree treebuilder.Tree
}

func Compile(projectPath string, opts Options) error {
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
	_, err = parseProjectFiles(projectPath, absEntrypointPath, diagnosticsBag, opts)
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
func parseProjectFiles(projectRootPath string, absEntrypointPath string, diagnosticsBag *diagnostics.Bag, opts Options) (*parseProjectFilesResult, error) {
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
		res, err := parseFile(projectRootPath, pathsToParse[0], src, diagnosticsBag.For(diagnosticFile), opts)
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
func parseFile(projectRootPath string, filePath string, src []byte, reporter *diagnostics.Reporter, opts Options) (*parseFileResult, error) {
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

	EmitDebugArtifact(opts.EmitTokens, "tokens", essentialTokens.String, projectRootPath, opts.Out, filePath)

	startAst := time.Now()
	tree := treebuilder.BuildTree(filePath, essentialTokens.Tokens, reporter)
	astTime := time.Since(startAst)
	performance.ReportHeapSize("ast")

	EmitDebugArtifact(opts.EmitTree, "tree", tree.String, projectRootPath, opts.Out, filePath)

	fmt.Printf("Lexer: %s, ASI: %s, AST: %s\n", lexerTime, asiTime, astTime)

	return &parseFileResult{
		Imports: []string{},
		Tree: tree,
	}, nil
}

func EmitDebugArtifact(emitFlagValue bool, typeName string, dataDumpFunc func() (string, error), projectRootPath string, outRootPath string, compiledFilePath string) error {
	if !emitFlagValue {
		return nil
	}

	filePathRelativeToProjectRoot, err := filepath.Rel(projectRootPath, compiledFilePath)
	if err != nil {
		return err
	}

	fullArtifactPath := path.Join(outRootPath, "debug", typeName, filePathRelativeToProjectRoot + ".json")
	fileParentDir := path.Dir(fullArtifactPath)

	if err := os.MkdirAll(fileParentDir, os.ModePerm); err != nil {
		return err
	}

	data, err := dataDumpFunc()
	if err != nil {
		return err
	}

	if err := os.WriteFile(fullArtifactPath, []byte(data), 0644); err != nil {
		return err
	}

	return nil
}

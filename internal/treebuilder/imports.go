package treebuilder

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"martinpetr.dev/kina/compiler/internal/utils"
	"martinpetr.dev/kina/compiler/projectConfig"
)

type ResolvedImportMember struct {
	Name  string
	Alias string
}

type ResolvedImport struct {
	ResolvedPath string
	Members 	[]ResolvedImportMember
}

type ResolvedImports struct {
	imports []ResolvedImport
}

func (t *Tree) ResolveImports(projectRoot string, currentFilePath string) ResolvedImports {
	var fileNode = t.Node
	var imports []ResolvedImport = make([]ResolvedImport, 0)

	for _, child := range fileNode.Children {
		switch child := child.(type) {
			case importNode:
				resolvedPath, err := resolveImportPath(projectRoot, currentFilePath, child.ModuleName)
				if err != nil {
					unquoted, _ := utils.GetStringLiteralValue(child.ModuleName)

					// TODO: Replace with diagnostics
					fmt.Printf("Error resolving import path for module '%s': %v\n", unquoted, err)
					continue
				}

				var members []ResolvedImportMember = make([]ResolvedImportMember, 0)
				for _, member := range child.Members {
					members = append(members, ResolvedImportMember{
						Name:  member.Name,
						Alias: member.Alias,
					})
				}

				imports = append(imports, ResolvedImport{
					ResolvedPath: resolvedPath,
					Members: members,
				})
			default:
				continue
		}
	}

	return ResolvedImports{
		imports: imports,
	}
}

func (i *ResolvedImports) GetPaths() []string {
	paths := make([]string, 0, len(i.imports))

	for _, imp := range i.imports {
		// Check if the resolved path is already in the list
		if !slices.Contains(paths, imp.ResolvedPath) {
			paths = append(paths, imp.ResolvedPath)
		}
	}

	return paths
}

// Use module resolving strategy to try to find the correct path for the module name.
// - ./..., ../... are resolved relative to the current file path. (Absolute paths are not allowed)
// - everything else is resolved as a module under its own stragegy:
//   - must be in format author.package with optional path (author.package.path.to.module)
//   - resolved as <project root>/.kina_modules/AUTHOR/PACKAGE and there the entry file (usually src/lib.kin)
//   - when path is provided, it is relative to the entry of the module (e.g. src/lib.kin) and must point to a file (lastSegment.kin) or a directory (lastSegment/lib.kin)
func resolveImportPath(projectRoot, currentFilePath, moduleName string) (string, error) {
	unescapedModuleName, err := utils.GetStringLiteralValue(moduleName)
	if err != nil {
		return "", err
	}

	// Relative path resolution
	if strings.HasPrefix(unescapedModuleName, "./") || strings.HasPrefix(unescapedModuleName, "../") {
		resolvedPath := filepath.Join(filepath.Dir(currentFilePath), unescapedModuleName)
		absResolvedPath, err := filepath.Abs(resolvedPath)
		if err != nil {
			return "", err
		}

		return absResolvedPath, nil
	}

	if !strings.Contains(unescapedModuleName, ".") {
		return "", fmt.Errorf("Invalid module name '%s'. Must be in format 'author.package' with optional path (dot separated)", unescapedModuleName)
	}

	author, packageName, modulePath := utils.SplitLibraryModuleName(unescapedModuleName)

	kinaModulesRoot := path.Join(projectRoot, ".kina_modules")
	moduleRootDir := path.Join(kinaModulesRoot, author, packageName)
	moduleConfigPath := path.Join(moduleRootDir, "kina.toml")

	// Check if module config exists
	if _, err := filepath.Abs(moduleConfigPath); err != nil {
		return "", fmt.Errorf("Module '%s.%s' not found at '%s'", author, packageName, moduleRootDir)
	}

	// Read the module config to get the entry file
	f, err := os.Open(moduleConfigPath)
	if err != nil {
		return "", fmt.Errorf("Failed to open module config at '%s': %v", moduleConfigPath, err)
	}
	defer f.Close()

	moduleConfig, err := projectConfig.ParseFile(moduleConfigPath)
	if err != nil {
		return "", fmt.Errorf("Failed to parse module config at '%s': %v", moduleConfigPath, err)
	}

	// If there is no path provided, return the entry file of the module
	if modulePath == "" {
		entryFilePath := path.Join(moduleRootDir, moduleConfig.Project.Entry)
		return entryFilePath, nil
	}

	// Get path to the dir of the entry file of the module
	entryFileDir := filepath.Dir(moduleConfig.Project.Entry)

	fileKinPath := path.Join(moduleRootDir, entryFileDir, modulePath + ".kin")
	dirKinPath := path.Join(moduleRootDir, entryFileDir, modulePath, "lib.kin")

	// Check if the file exists
	if _, err := os.Stat(fileKinPath); err == nil {
		return fileKinPath, nil
	}

	// Check if the directory exists
	if _, err := os.Stat(dirKinPath); err == nil {
		return dirKinPath, nil
	}

	return "", fmt.Errorf("Module '%s.%s' does not have a file or directory at path '%s'", author, packageName, modulePath)
}

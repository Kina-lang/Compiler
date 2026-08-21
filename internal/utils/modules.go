package utils

import "strings"

func SplitLibraryModuleName(moduleName string) (string, string, string) {
	author := ""
	pkg := ""
	pathSegments := []string{}

	for i, part := range strings.Split(moduleName, ".") {
		switch i {
		case 0:
			author = part
		case 1:
			pkg = part
		default:
			pathSegments = append(pathSegments, part)
		}
	}

	return author, pkg, strings.Join(pathSegments, "/")
}

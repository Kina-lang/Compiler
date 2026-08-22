package tools

import "os/exec"

func ClangCommand() string {
	return "clang" // TODO: resolve to toolchain installation dir
}

func ClangBuildObject(target string, sourceCode []byte, outputPath string) {
	cmd := exec.Command(ClangCommand(), "-target", target, "-x", "ir", "-c", "-", "-o", outputPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	_, err = stdin.Write(sourceCode)
	if err != nil {
		panic(err)
	}

	err = stdin.Close()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func ClangLink(target string, includedFilePaths []string, outputPath string) {
	cmd := exec.Command(ClangCommand(), "-target", target, "-o", outputPath, "-fuse-ld=mold")
	cmd.Args = append(cmd.Args, includedFilePaths...)

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

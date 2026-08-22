package tools

import "os/exec"

func LlcCommand() string {
	return "llc" // TODO: resolve to toolchain installation dir
}

func LlcBuildObject(sourceCode []byte, outputPath string) {
	cmd := exec.Command(LlcCommand(), "-filetype=obj", "-O0", "-o", outputPath, "-")

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

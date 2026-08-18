package main

import "martinpetr.dev/kina/internal/driver"

func main() {
	driver.Compile("src/main.kin", driver.Options{
		Out: "build/main",
	});
}

package compiler

import "fmt"

func Greet(name string) {
	fmt.Println(GetGreeting(name))
}

func GetGreeting(name string) string {
	return "Hello, "+name+"!"
}

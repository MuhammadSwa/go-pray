package main

import "fmt"

// show help message
func help() {
	msg := `Valid arguments:
	list - List all prayers
	next - Show time left for next prayer
	date - Show hijri date 
	help - Show help message `
	fmt.Println(msg)
}

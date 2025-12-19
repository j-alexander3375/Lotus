package main

import (
	"fmt"
	"os"
)

func printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func fprintf(w *os.File, format string, a ...interface{}) {
	fmt.Fprintf(w, format, a...)
}

func println(a ...interface{}) {
	fmt.Println(a...)
}

func sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}

func sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}
func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func fatalln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
func logf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func logln(a ...interface{}) {
	fmt.Println(a...)
}

package main

import "fmt"

func Assert(cond bool, msg string) {
	if !cond {
		panic(fmt.Sprintf("ASSERTION FAILED: %s", msg))
	}
}

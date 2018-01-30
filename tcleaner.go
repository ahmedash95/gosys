package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearCli() {
	out := runtime.GOOS
	if runtime.GOOS == "darwin" {
		out = "linux"
	}
	value, ok := clear[out] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                 //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		fmt.Println(runtime.GOOS)
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

package main

import (
	"fmt"
	"os"

	"wolf404/compiler/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		cmd.InitProject(os.Args[2:])
	case "run":
		cmd.RunFile(os.Args[2:])
	case "build":
		cmd.BuildFile(os.Args[2:])
	case "pack":
		cmd.ManagePackage(os.Args[2:])
	case "test":
		cmd.RunTests(os.Args[2:])
	case "fmt":
		cmd.FormatFile(os.Args[2:])
	case "version":
		cmd.ShowVersion()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("🐺 Wolf404 Compiler v0.1.0")
	fmt.Println("\nUsage:")
	fmt.Println("  wlf init <project-name>       Create new Wolf404 project")
	fmt.Println("  wlf run <file.wlf>            Run Wolf404 file in dev mode")
	fmt.Println("  wlf build <file.wlf>          Build executable binary")
	fmt.Println("  wlf pack add <package>        Add package to project")
	fmt.Println("  wlf test                      Run tests")
	fmt.Println("  wlf fmt <file.wlf>            Format Wolf404 code")
	fmt.Println("  wlf version                   Show version")
	fmt.Println("\nExamples:")
	fmt.Println("  wlf init my-web-app")
	fmt.Println("  wlf run main.wlf")
	fmt.Println("  wlf build main.wlf -o app")
}

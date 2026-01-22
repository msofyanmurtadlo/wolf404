package main

import (
	"fmt"
	"os"

	"wolf404/compiler/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.RunREPL()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		cmd.InitProject(os.Args[2:])
	case "gas":
		cmd.RunFile(os.Args[2:])
	case "gawe:model":
		cmd.MakeModel(os.Args[2:])
	case "gawe:controller":
		cmd.MakeController(os.Args[2:])
	case "gawe:middleware":
		cmd.MakeMiddleware(os.Args[2:])
	case "migrate":
		cmd.RunMigrations(os.Args[2:])
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
	fmt.Println("Created by ishowpen")
	fmt.Println("\nUsage:")
	fmt.Println("  wlf init <project-name>       Create new Wolf404 project")
	fmt.Println("  wlf gas <file.wlf>            Run Wolf404 file in dev mode")
	fmt.Println("  wlf gas server                Start the Wolf404 server")
	fmt.Println("  wlf gawe:model <Name>         Generate Model and Migration")
	fmt.Println("  wlf gawe:controller <Name>    Generate Controller")
	fmt.Println("  wlf gawe:middleware <Name>    Generate Middleware")
	fmt.Println("  wlf migrate                   Run database migrations")
	fmt.Println("  wlf build <file.wlf>          Build executable binary")
	fmt.Println("  wlf test                      Run tests")
	fmt.Println("  wlf fmt <file.wlf>            Format Wolf404 code")
	fmt.Println("  wlf version                   Show version")
	fmt.Println("\nExamples:")
	fmt.Println("  wlf init my-web-app")
	fmt.Println("  wlf gas server")
	fmt.Println("  wlf gawe:model Product")
	fmt.Println("  wlf gawe:controller ProductController")
	fmt.Println("  wlf migrate")
}

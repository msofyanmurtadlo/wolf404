package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"wolf404/compiler/evaluator"
	"wolf404/compiler/lexer"
	"wolf404/compiler/object"
	"wolf404/compiler/parser"
	"wolf404/compiler/repl"
)

func RunREPL() {
	fmt.Println("🐺 Wolf404 v1.1 Interactive Shell")
	fmt.Println("Type commands to evaluate them.")
	repl.Start(os.Stdin, os.Stdout)
}

func InitProject(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: wlf init <project_name>")
		return
	}
	projectName := args[0]
	fmt.Printf("Initializing Wolf404 project: %s...\n", projectName)
	
	// Create directories
	dirs := []string{
		projectName,
		filepath.Join(projectName, "packs"),
		filepath.Join(projectName, "views"),
		filepath.Join(projectName, "public"),
		filepath.Join(projectName, "bin"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			return
		}
	}

	// Create main.wlf
	mainContent := `summon web_server
summon router

hunt main()
    howl("🐺 Wolf404 Server Starting...")
    web_server.start(8080)
`
	ioutil.WriteFile(filepath.Join(projectName, "main.wlf"), []byte(mainContent), 0644)
	fmt.Println("Project initialized successfully! 🐺")
}

func RunFile(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: wlf run <file.wlf>")
		return
	}
	filename := args[0]
	fmt.Printf("Running %s...\n", filename)
	
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Parse
	fmt.Println("Starting Lexer/Parser...")
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	fmt.Printf("Parsed %d statements\n", len(program.Statements))

	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		return
	}

	// Evaluate
	fmt.Println("\n--- AST Structure ---")
	fmt.Println(program.String())

	fmt.Println("\n--- Execution Output ---")
	
	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)
	
	if evaluated != nil {
		// Only print result if it's an error or significant?
		// Usually programs print via 'howl', so returning inspection is strictly for REPL.
		if evaluated.Type() == object.ERROR_OBJ {
			fmt.Println(evaluated.Inspect())
		}
	}
}

func BuildFile(args []string) {
	fmt.Println("Build not implemented yet")
}

func ManagePackage(args []string) {
	fmt.Println("Package manager not implemented yet")
}

func RunTests(args []string) {
	fmt.Println("Tests not implemented yet")
}

func FormatFile(args []string) {
	fmt.Println("Formatter not implemented yet")
}

func ShowVersion() {
	fmt.Println("🐺 Wolf404 Compiler v0.1.0")
}

func printParserErrors(errors []string) {
	fmt.Println("🐺 Wolf404 encountered errors during parsing:")
	for _, msg := range errors {
		fmt.Printf("\t%s\n", msg)
	}
}

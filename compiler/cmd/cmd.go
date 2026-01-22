package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
	"wolf404/compiler/evaluator"
	"wolf404/compiler/lexer"
	"wolf404/compiler/object"
	"wolf404/compiler/parser"
	"wolf404/compiler/repl"

	_ "modernc.org/sqlite"
)

func RunREPL() {
	fmt.Println("🐺 Wolf404 v1.1 Interactive Shell")
	fmt.Println("Created by ishowpen")
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
		fmt.Println("Usage: wlf gas <file.wlf> or wlf gas server")
		return
	}
	filename := args[0]

	// Artisan-like shortcut: "wlf gas server" -> runs "server.wlf"
	if filename == "server" {
		if _, err := os.Stat("server.wlf"); err == nil {
			filename = "server.wlf"
		} else {
			fmt.Println("❌ Error: 'server.wlf' not found. Make sure you are in the project root.")
			return
		}
	}

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
	fmt.Println("Created by ishowpen")
}

func RunMigrations(args []string) {
	fmt.Println("🔄 Running migrations...")

	// Open database for tracking
	db, err := sql.Open("sqlite", "database/database.db")
	if err != nil {
		fmt.Printf("❌ Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()

	// Create migrations table if not exists
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS migrations (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, run_at DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		fmt.Printf("❌ Error creating migrations tracker: %v\n", err)
		return
	}

	migrationsDir := "database/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		fmt.Printf("❌ Error reading migrations directory: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("ℹ️  No migrations found")
		return
	}

	runCount := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".wlf" {
			// Check if already run
			var id int
			err = db.QueryRow("SELECT id FROM migrations WHERE name = ?", file.Name()).Scan(&id)
			if err == nil {
				// Already run, skip
				continue
			}

			migrationPath := filepath.Join(migrationsDir, file.Name())
			fmt.Printf("⚡ Running: %s\n", file.Name())

			// Run the migration file
			content, err := ioutil.ReadFile(migrationPath)
			if err != nil {
				fmt.Printf("   ❌ Error reading: %v\n", err)
				continue
			}

			l := lexer.New(string(content))
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				fmt.Printf("   ❌ Parse errors in %s\n", file.Name())
				continue
			}

			env := object.NewEnvironment()
			evaluator.Eval(program, env)

			// Record in tracker
			_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", file.Name())
			if err != nil {
				fmt.Printf("   ❌ Error recording migration: %v\n", err)
			} else {
				fmt.Printf("   ✅ Migrated: %s\n", file.Name())
				runCount++
			}
		}
	}

	if runCount == 0 {
		fmt.Println("ℹ️  Nothing to migrate.")
	} else {
		fmt.Printf("\n🐺 %d migrations completed!\n", runCount)
	}
}

func MakeController(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: wlf gawe:controller <ControllerName>")
		fmt.Println("Example: wlf gawe:controller ProductController")
		return
	}

	controllerName := args[0]
	if !strings.HasSuffix(controllerName, "Controller") {
		controllerName += "Controller"
	}

	controllerPath := filepath.Join("app", "Controllers", controllerName+".wlf")
	controllerContent := fmt.Sprintf(`// app/Controllers/%s.wlf
// Created by ishowpen

gerombolan %s
    garap index($req)
        balekno http_json({"message": "Index method"})
    
    garap show($req, $id)
        balekno http_json({"id": $id})
    
    garap store($req)
        balekno http_json({"message": "Created"})
    
    garap update($req, $id)
        balekno http_json({"message": "Updated"})
    
    garap destroy($req, $id)
        balekno http_json({"message": "Deleted"})
`, controllerName, controllerName)

	err := os.MkdirAll(filepath.Dir(controllerPath), 0755)
	if err != nil {
		fmt.Printf("❌ Error creating directory: %v\n", err)
		return
	}

	err = ioutil.WriteFile(controllerPath, []byte(controllerContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error creating controller: %v\n", err)
		return
	}

	fmt.Printf("✅ Controller created: %s\n", controllerPath)
}

func MakeMiddleware(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: wlf gawe:middleware <MiddlewareName>")
		fmt.Println("Example: wlf gawe:middleware CheckAdmin")
		return
	}

	middlewareName := args[0]
	middlewarePath := filepath.Join("app", "Middleware", middlewareName+".wlf")
	middlewareContent := fmt.Sprintf(`// app/Middleware/%s.wlf
// Created by ishowpen

gerombolan %s
    garap handle($request)
        // Add your middleware logic here
        balekno {"error": salah}
`, middlewareName, middlewareName)

	err := os.MkdirAll(filepath.Dir(middlewarePath), 0755)
	if err != nil {
		fmt.Printf("❌ Error creating directory: %v\n", err)
		return
	}

	err = ioutil.WriteFile(middlewarePath, []byte(middlewareContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error creating middleware: %v\n", err)
		return
	}

	fmt.Printf("✅ Middleware created: %s\n", middlewarePath)
}

func MakeModel(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: wlf gawe:model <ModelName>")
		fmt.Println("Example: wlf gawe:model Product")
		return
	}

	modelName := args[0]

	// Capitalize first letter
	if len(modelName) > 0 {
		modelName = strings.ToUpper(string(modelName[0])) + modelName[1:]
	}

	// Create Model file
	modelPath := filepath.Join("app", "Models", modelName+".wlf")
	modelContent := fmt.Sprintf(`// app/Models/%s.wlf
// Created by ishowpen

gerombolan %s
    garap init($db)
        $this.db = $db
        $this.table = "%s"
    
    garap all()
        $query = "SELECT * FROM " + $this.table
        balekno db_query($query)
    
    garap find($id)
        $query = "SELECT * FROM " + $this.table + " WHERE id = ?"
        balekno db_query($query)
    
    garap create($data)
        // Implementation for creating record
        balekno bener
    
    garap update($id, $data)
        // Implementation for updating record
        balekno bener
    
    garap delete($id)
        $query = "DELETE FROM " + $this.table + " WHERE id = ?"
        balekno db_exec($query)
`, modelName, modelName, strings.ToLower(modelName)+"s")

	err := os.MkdirAll(filepath.Dir(modelPath), 0755)
	if err != nil {
		fmt.Printf("❌ Error creating directory: %v\n", err)
		return
	}

	err = ioutil.WriteFile(modelPath, []byte(modelContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error creating model: %v\n", err)
		return
	}

	fmt.Printf("✅ Model created: %s\n", modelPath)

	// Create Migration file
	timestamp := time.Now().Format("20060102_150405")
	migrationName := fmt.Sprintf("%s_create_%s_table", timestamp, strings.ToLower(modelName)+"s")
	migrationPath := filepath.Join("database", "migrations", migrationName+".wlf")

	migrationContent := fmt.Sprintf(`// database/migrations/%s.wlf
// Created by ishowpen

garap up($db)
    $sql = "CREATE TABLE IF NOT EXISTS %s (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )"
    db_exec($sql)
    ketok("✅ Table %s created")

garap down($db)
    $sql = "DROP TABLE IF EXISTS %s"
    db_exec($sql)
    ketok("✅ Table %s dropped")
`, migrationName, strings.ToLower(modelName)+"s", strings.ToLower(modelName)+"s", strings.ToLower(modelName)+"s", strings.ToLower(modelName)+"s")

	err = os.MkdirAll(filepath.Dir(migrationPath), 0755)
	if err != nil {
		fmt.Printf("❌ Error creating migrations directory: %v\n", err)
		return
	}

	err = ioutil.WriteFile(migrationPath, []byte(migrationContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error creating migration: %v\n", err)
		return
	}

	fmt.Printf("✅ Migration created: %s\n", migrationPath)
	fmt.Println("\n🐺 Model and Migration created successfully!")
}

func printParserErrors(errors []string) {
	fmt.Println("🐺 Wolf404 encountered errors during parsing:")
	for _, msg := range errors {
		fmt.Printf("\t%s\n", msg)
	}
}

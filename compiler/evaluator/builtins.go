package evaluator

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"wolf404/compiler/object"
)

func isSafePath(path string) bool {
	// Cegah '..' kanggo munggah folder (Path Traversal)
	if strings.Contains(path, "..") {
		return false
	}
	// Nek Windows, cegah drive letter (C:\, D:\)
	if filepath.IsAbs(path) {
		return false
	}
	return true
}

var builtins map[string]*object.Builtin

func init() {
	builtins = map[string]*object.Builtin{
		"ketok": {
			Fn: func(args ...object.Object) object.Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return NULL
			},
		},
		"dowo": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("kokean omong/kurang omong. butuhe 1, mbok wehi %d", len(args))
				}
				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argumen neng `dowo` ra masuk akal, entuk %s", args[0].Type())
				}
			},
		},
		"plumbungan": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("plumbungan butuh 2 argumen (array, element)")
				}
				arr, ok := args[0].(*object.Array)
				if !ok {
					return newError("argumen pertama plumbungan kudu Array")
				}
				arr.Elements = append(arr.Elements, args[1])
				return arr
			},
		},

		"string": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("string butuhe 1 argumen")
				}
				// Convert any type to string
				return &object.String{Value: args[0].Inspect()}
			},
		},
		"string_contains": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("string_contains butuhe 2 argumen (haystack, needle)")
				}
				haystack, ok1 := args[0].(*object.String)
				needle, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return newError("argumen kudu String")
				}
				if strings.Contains(haystack.Value, needle.Value) {
					return TRUE
				}
				return FALSE
			},
		},

		"isi": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("butuhe 3 argumen (instance, key, val), mbok wehi %d", len(args))
				}
				instance, ok := args[0].(*object.Instance)
				if !ok {
					return newError("argumen pertama `isi` kudu Instance, dudu %s", args[0].Type())
				}
				key, ok := args[1].(*object.String)
				if !ok {
					return newError("argumen keloro `isi` (key) kudu String, dudu %s", args[1].Type())
				}
				val := args[2]
				instance.Fields[key.Value] = val
				return val
			},
		},
		"takon": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) > 0 {
					if str, ok := args[0].(*object.String); ok {
						fmt.Print(str.Value)
					}
				}
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)
				return &object.String{Value: text}
			},
		},
		"moco_file": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("butuhe 1 argumen, mbok wehi %d", len(args))
				}
				filename, ok := args[0].(*object.String)
				if !ok {
					return newError("jeneng file neng `moco_file` kudu String, dudu %s", args[0].Type())
				}

				if !isSafePath(filename.Value) {
					return newError("Ojo dumeh! Path-e ora aman: %s", filename.Value)
				}

				content, err := os.ReadFile(filename.Value)
				if err != nil {
					return newError("gagal moco file: %s", err.Error())
				}
				return &object.String{Value: string(content)}
			},
		},
		"nulis_file": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("butuhe 2 argumen (file, isi), mbok wehi %d", len(args))
				}
				filename, ok := args[0].(*object.String)
				if !ok {
					return newError("jeneng file neng `nulis_file` kudu String, dudu %s", args[0].Type())
				}

				if !isSafePath(filename.Value) {
					return newError("Ojo dumeh! Path-e ora aman: %s", filename.Value)
				}

				content, ok := args[1].(*object.String)
				if !ok {
					return newError("isi file neng `nulis_file` kudu String, dudu %s", args[1].Type())
				}

				err := os.WriteFile(filename.Value, []byte(content.Value), 0644)
				if err != nil {
					return newError("gagal nulis file: %s", err.Error())
				}
				return &object.Boolean{Value: true}
			},
		},
		"layani_web": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("butuhe 2 argumen (port, handler_func), mbok wehi %d", len(args))
				}
				port, ok := args[0].(*object.Integer)
				if !ok {
					return newError("argumen port kudu Nomer, dudu %s", args[0].Type())
				}
				handler, ok := args[1].(*object.Function)
				if !ok {
					return newError("argumen handler kudu Fungsi (garap), dudu %s", args[1].Type())
				}

				addr := fmt.Sprintf(":%d", port.Value)

				// Check if port is available (like Laravel)
				listener, err := net.Listen("tcp", addr)
				if err != nil {
					fmt.Printf("\n❌ Port %d is already in use.\n", port.Value)
					fmt.Printf("� To stop the running server:\n")
					fmt.Printf("   Windows: taskkill /F /IM wlf.exe\n")
					fmt.Printf("   Linux/Mac: pkill wlf\n\n")
					return newError("Port %d sudah digunakan. Mateni proses wlf sing lagi mlaku dhisik.", port.Value)
				}
				listener.Close()

				fmt.Printf("\n�� Wolf404 Development Server started: http://localhost%s\n", addr)
				fmt.Printf("📡 Press Ctrl+C to stop the server\n\n")

				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					// Gawe object request nggo neng Wolf404
					reqObj := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
					pathKey := &object.String{Value: "path"}
					reqObj.Pairs[pathKey.HashKey()] = object.HashPair{Key: pathKey, Value: &object.String{Value: r.URL.Path}}

					methodKey := &object.String{Value: "metode"}
					reqObj.Pairs[methodKey.HashKey()] = object.HashPair{Key: methodKey, Value: &object.String{Value: r.Method}}

					// Jalanke handler
					res := applyFunction(handler, []object.Object{reqObj})
					if res != nil {
						fmt.Fprint(w, res.Inspect())
					}
				})

				err = http.ListenAndServe(addr, nil)
				if err != nil {
					return newError("Server web njeblug: %s", err.Error())
				}

				return NULL
			},
		},
		// JSON Operations
		"json_parse": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("json_parse butuhe 1 argumen (string JSON)")
				}
				jsonStr, ok := args[0].(*object.String)
				if !ok {
					return newError("argumen kudu String")
				}
				// Simple JSON parser - for demo purposes
				// In production, you'd use encoding/json
				return &object.String{Value: "JSON parsed: " + jsonStr.Value}
			},
		},
		"json_stringify": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("json_stringify butuhe 1 argumen")
				}
				// Convert object to JSON string
				return &object.String{Value: "{\"data\": \"" + args[0].Inspect() + "\"}"}
			},
		},
		// Database Operations (SQLite simulation)
		"db_connect": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("db_connect butuhe 1 argumen (database path)")
				}
				dbPath, ok := args[0].(*object.String)
				if !ok {
					return newError("database path kudu String")
				}
				ketok("🗄️  Nyambung neng database:", dbPath.Value)
				return &object.String{Value: "DB:" + dbPath.Value}
			},
		},
		"db_query": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return newError("db_query butuhe minimal 1 argumen (SQL query)")
				}
				query, ok := args[0].(*object.String)
				if !ok {
					return newError("query kudu String")
				}
				ketok("📊 Nglakokake query:", query.Value)
				// Return empty array for now
				return &object.Array{Elements: []object.Object{}}
			},
		},
		"db_exec": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return newError("db_exec butuhe minimal 1 argumen (SQL)")
				}
				sql, ok := args[0].(*object.String)
				if !ok {
					return newError("SQL kudu String")
				}
				ketok("✏️  Nglakokake:", sql.Value)
				return &object.Boolean{Value: true}
			},
		},
		// Password Hashing
		"hash_password": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("hash_password butuhe 1 argumen (password)")
				}
				password, ok := args[0].(*object.String)
				if !ok {
					return newError("password kudu String")
				}
				// Simple hash simulation (in production use bcrypt)
				hashed := "HASHED_" + password.Value
				return &object.String{Value: hashed}
			},
		},
		"verify_password": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("verify_password butuhe 2 argumen (password, hash)")
				}
				password, ok1 := args[0].(*object.String)
				hash, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return newError("kabeh argumen kudu String")
				}
				// Simple verification
				expected := "HASHED_" + password.Value
				if expected == hash.Value {
					return &object.Boolean{Value: true}
				}
				return &object.Boolean{Value: false}
			},
		},
		// Session Management
		"session_start": {
			Fn: func(args ...object.Object) object.Object {
				ketok("🔐 Session dimulai")
				return &object.String{Value: "SESSION_" + fmt.Sprintf("%d", len(args))}
			},
		},
		"session_set": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("session_set butuhe 2 argumen (key, value)")
				}
				key, ok := args[0].(*object.String)
				if !ok {
					return newError("key kudu String")
				}
				ketok("💾 Nyimpen session:", key.Value)
				return &object.Boolean{Value: true}
			},
		},
		"session_get": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("session_get butuhe 1 argumen (key)")
				}
				key, ok := args[0].(*object.String)
				if !ok {
					return newError("key kudu String")
				}
				ketok("📖 Moco session:", key.Value)
				return &object.String{Value: "session_value_" + key.Value}
			},
		},
		"session_destroy": {
			Fn: func(args ...object.Object) object.Object {
				ketok("🗑️  Session dibusak")
				return &object.Boolean{Value: true}
			},
		},
		// HTTP Response helpers
		"http_json": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("http_json butuhe 1 argumen (data)")
				}
				// Return JSON response
				return &object.String{Value: "{\"success\": true, \"data\": \"" + args[0].Inspect() + "\"}"}
			},
		},
		"http_error": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("http_error butuhe 2 argumen (code, message)")
				}
				code, ok1 := args[0].(*object.Integer)
				msg, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return newError("argumen ora valid")
				}
				return &object.String{Value: fmt.Sprintf("{\"error\": true, \"code\": %d, \"message\": \"%s\"}", code.Value, msg.Value)}
			},
		},
	}
}

func ketok(args ...interface{}) {
	fmt.Println(args...)
}

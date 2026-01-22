package evaluator

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"wolf404/compiler/object"

	_ "modernc.org/sqlite"
)

func isSafePath(path string) bool {
	if strings.Contains(path, "..") {
		return false
	}
	if filepath.IsAbs(path) {
		return false
	}
	return true
}

var dbConnections = make(map[string]*sql.DB)
var sessions = make(map[string]object.Object)

func convertToWolfObject(v interface{}) object.Object {
	switch val := v.(type) {
	case int64:
		return &object.Integer{Value: val}
	case float64:
		return &object.String{Value: fmt.Sprintf("%f", val)}
	case string:
		return &object.String{Value: val}
	case []byte:
		return &object.String{Value: string(val)}
	case bool:
		return &object.Boolean{Value: val}
	case map[string]interface{}:
		hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
		for k, v := range val {
			key := &object.String{Value: k}
			hash.Pairs[key.HashKey()] = object.HashPair{Key: key, Value: convertToWolfObject(v)}
		}
		return hash
	case []interface{}:
		elements := make([]object.Object, len(val))
		for i, v := range val {
			elements[i] = convertToWolfObject(v)
		}
		return &object.Array{Elements: elements}
	case nil:
		return NULL
	default:
		return &object.String{Value: fmt.Sprintf("%v", val)}
	}
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
					return newError("kokean omong/kurang omong. butuhe 1")
				}
				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				case *object.Hash:
					return &object.Integer{Value: int64(len(arg.Pairs))}
				default:
					return newError("argumen neng `dowo` ra masuk akal")
				}
			},
		},
		"plumbungan": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("plumbungan butuh 2 argumen")
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
				return &object.String{Value: args[0].Inspect()}
			},
		},
		"string_contains": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("string_contains butuhe 2 argumen")
				}
				return &object.Boolean{Value: strings.Contains(args[0].Inspect(), args[1].Inspect())}
			},
		},
		"string_split": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("string_split butuhe 2 argumen")
				}
				parts := strings.Split(args[0].Inspect(), args[1].Inspect())
				elements := make([]object.Object, len(parts))
				for i, p := range parts {
					elements[i] = &object.String{Value: p}
				}
				return &object.Array{Elements: elements}
			},
		},
		"string_replace": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return newError("string_replace butuhe 3 argumen")
				}
				return &object.String{Value: strings.ReplaceAll(args[0].Inspect(), args[1].Inspect(), args[2].Inspect())}
			},
		},
		"string_regex_match": {
			Fn: func(args ...object.Object) object.Object {
				match, _ := regexp.MatchString(args[1].Inspect(), args[0].Inspect())
				return &object.Boolean{Value: match}
			},
		},
		"string_regex_capture": {
			Fn: func(args ...object.Object) object.Object {
				re, _ := regexp.Compile(args[1].Inspect())
				matches := re.FindStringSubmatch(args[0].Inspect())
				if len(matches) == 0 {
					return &object.Array{}
				}
				elements := make([]object.Object, len(matches)-1)
				for i := 1; i < len(matches); i++ {
					elements[i-1] = &object.String{Value: matches[i]}
				}
				return &object.Array{Elements: elements}
			},
		},
		"html_escape": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("html_escape butuhe 1 argumen")
				}
				return &object.String{Value: html.EscapeString(args[0].Inspect())}
			},
		},
		"keys": {
			Fn: func(args ...object.Object) object.Object {
				hash, ok := args[0].(*object.Hash)
				if !ok {
					return newError("kudu Hash")
				}
				elements := make([]object.Object, 0, len(hash.Pairs))
				for _, pair := range hash.Pairs {
					elements = append(elements, pair.Key)
				}
				return &object.Array{Elements: elements}
			},
		},
		"moco_file": {
			Fn: func(args ...object.Object) object.Object {
				path := args[0].Inspect()
				if !isSafePath(path) {
					return newError("path bahaya")
				}
				content, err := os.ReadFile(path)
				if err != nil {
					return newError("gagal moco")
				}
				return &object.String{Value: string(content)}
			},
		},
		"layani_web": {
			Fn: func(args ...object.Object) object.Object {
				port := args[0].(*object.Integer).Value
				handler := args[1].(*object.Function)
				addr := fmt.Sprintf(":%d", port)

				fmt.Printf("\n🐺 Wolf404 Development Server started: http://localhost%s\n", addr)
				http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
					// 1. Base Request Info
					reqData := make(map[string]interface{})
					reqData["path"] = r.URL.Path
					reqData["metode"] = r.Method

					// 2. Parse Query Params
					query := make(map[string]interface{})
					for k, v := range r.URL.Query() {
						if len(v) == 1 {
							query[k] = v[0]
						} else {
							query[k] = v
						}
					}
					reqData["query"] = query

					// 3. Parse Body (JSON or Form)
					body := make(map[string]interface{})
					if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
						contentType := r.Header.Get("Content-Type")
						if strings.Contains(contentType, "application/json") {
							json.NewDecoder(r.Body).Decode(&body)
						} else {
							r.ParseForm()
							for k, v := range r.PostForm {
								if len(v) == 1 {
									body[k] = v[0]
								} else {
									body[k] = v
								}
							}
						}
					}
					reqData["input"] = body

					// Convert to Wolf404 Object
					reqObj := convertToWolfObject(reqData)

					res := applyFunction(handler, []object.Object{reqObj})
					if res != nil {
						fmt.Fprint(w, res.Inspect())
					}
				})
				http.ListenAndServe(addr, nil)
				return NULL
			},
		},
		"db_connect": {
			Fn: func(args ...object.Object) object.Object {
				path := args[0].Inspect()
				db, err := sql.Open("sqlite", path)
				if err != nil {
					return newError("db error")
				}
				dbConnections["default"] = db
				return &object.String{Value: "connected"}
			},
		},
		"db_query": {
			Fn: func(args ...object.Object) object.Object {
				db := dbConnections["default"]
				if db == nil {
					return newError("no db")
				}
				var params []interface{}
				if len(args) > 1 {
					arr := args[1].(*object.Array)
					for _, el := range arr.Elements {
						params = append(params, el.Inspect())
					}
				}
				rows, err := db.Query(args[0].Inspect(), params...)
				if err != nil {
					return newError("query error")
				}
				defer rows.Close()
				cols, _ := rows.Columns()
				results := &object.Array{}
				for rows.Next() {
					vals := make([]interface{}, len(cols))
					ptr := make([]interface{}, len(cols))
					for i := range vals {
						ptr[i] = &vals[i]
					}
					rows.Scan(ptr...)
					row := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}
					for i, name := range cols {
						k := &object.String{Value: name}
						v := convertToWolfObject(vals[i])
						row.Pairs[k.HashKey()] = object.HashPair{Key: k, Value: v}
					}
					results.Elements = append(results.Elements, row)
				}
				return results
			},
		},
		"db_exec": {
			Fn: func(args ...object.Object) object.Object {
				db := dbConnections["default"]
				if db == nil {
					return newError("no db")
				}
				var params []interface{}
				if len(args) > 1 {
					arr := args[1].(*object.Array)
					for _, el := range arr.Elements {
						params = append(params, el.Inspect())
					}
				}
				_, err := db.Exec(args[0].Inspect(), params...)
				if err != nil {
					return newError("exec error")
				}
				return TRUE
			},
		},
		"hash_password": {
			Fn: func(args ...object.Object) object.Object {
				return &object.String{Value: "HASHED_" + args[0].Inspect()}
			},
		},
		"verify_password": {
			Fn: func(args ...object.Object) object.Object {
				return &object.Boolean{Value: "HASHED_"+args[0].Inspect() == args[1].Inspect()}
			},
		},
		"http_json": {
			Fn: func(args ...object.Object) object.Object { return args[0] },
		},
		"http_error": {
			Fn: func(args ...object.Object) object.Object {
				return &object.String{Value: "Error: " + args[1].Inspect()}
			},
		},
		"session_set": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("session_set butuhe 2 argumen")
				}
				key := args[0].Inspect()
				sessions[key] = args[1]
				return TRUE
			},
		},
		"session_get": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("session_get butuhe 1 argumen")
				}
				key := args[0].Inspect()
				val, ok := sessions[key]
				if !ok {
					return &object.String{Value: "kopong"}
				}
				return val
			},
		},
		"session_destroy": {
			Fn: func(args ...object.Object) object.Object {
				sessions = make(map[string]object.Object)
				return TRUE
			},
		},
	}
}

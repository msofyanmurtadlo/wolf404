package evaluator

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"wolf404/compiler/lexer"
	"wolf404/compiler/object"
	"wolf404/compiler/parser"

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
		"eval_wolf": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return newError("eval_wolf butuhe 1 argumen (code)")
				}
				code, ok := args[0].(*object.String)
				if !ok {
					return newError("code kudu String")
				}

				env := object.NewEnvironment()
				if len(args) > 1 {
					if hash, ok := args[1].(*object.Hash); ok {
						for _, pair := range hash.Pairs {
							if s, ok := pair.Key.(*object.String); ok {
								env.Set(s.Value, pair.Value)
							}
						}
					}
				}

				l := lexer.New(code.Value)
				p := parser.New(l)
				program := p.ParseProgram()
				if len(p.Errors()) != 0 {
					return newError("eval_wolf parse error: %s", p.Errors()[0])
				}

				return Eval(program, env)
			},
		},
		"render_template": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) < 1 {
					return newError("render_template butuhe 1 argumen (template)")
				}
				template, ok := args[0].(*object.String)
				if !ok {
					return newError("template kudu String")
				}

				data := make(map[string]object.Object)
				if len(args) > 1 {
					if hash, ok := args[1].(*object.Hash); ok {
						for _, pair := range hash.Pairs {
							if s, ok := pair.Key.(*object.String); ok {
								data[s.Value] = pair.Value
							}
						}
					}
				}

				raw := template.Value

				// 1. Resolve @leboke (Include)
				reInclude := regexp.MustCompile(`@leboke\s*\(\"(.*?)\"\)`)
				for reInclude.MatchString(raw) {
					raw = reInclude.ReplaceAllStringFunc(raw, func(m string) string {
						name := reInclude.FindStringSubmatch(m)[1]
						path := filepath.Join("resources", "views", strings.ReplaceAll(name, ".", "/")+".wlf")
						content, err := os.ReadFile(path)
						if err != nil {
							return "<!-- Error Include: " + path + " -->"
						}
						return string(content)
					})
				}

				// 2. STACKS logic (Extract @tumpuk before @warisan)
				stacks := make(map[string][]string)
				rePush := regexp.MustCompile(`(?s)@tumpuk\s*\(\"(.*?)\"\)(.*?)@punkyan_tumpuk`)
				matchesPush := rePush.FindAllStringSubmatch(raw, -1)
				for _, m := range matchesPush {
					name := m[1]
					content := m[2]
					stacks[name] = append(stacks[name], content)
				}
				// Remove push blocks from raw content
				raw = rePush.ReplaceAllString(raw, "")

				// 3. Resolve @warisan (Extend) and @bagean (Section)
				reExtend := regexp.MustCompile(`@warisan\s*\(\"(.*?)\"\)`)
				if reExtend.MatchString(raw) {
					layoutName := reExtend.FindStringSubmatch(raw)[1]
					layoutPath := filepath.Join("resources", "views", strings.ReplaceAll(layoutName, ".", "/")+".wlf")
					layoutContent, err := os.ReadFile(layoutPath)
					if err == nil {
						// Extract sections from child
						sections := make(map[string]string)
						reSection := regexp.MustCompile(`(?s)@bagean\s*\(\"(.*?)\"\)(.*?)@punkyan_bagean`)
						matches := reSection.FindAllStringSubmatch(raw, -1)
						for _, m := range matches {
							sections[m[1]] = m[2]
						}

						// Replace @panggonan (Yield) and @papan_tumpukan (Stack) in layout
						finalRaw := string(layoutContent)

						// Yields
						reYield := regexp.MustCompile(`@panggonan\s*\(\"(.*?)\"\)`)
						finalRaw = reYield.ReplaceAllStringFunc(finalRaw, func(m string) string {
							name := reYield.FindStringSubmatch(m)[1]
							if content, ok := sections[name]; ok {
								return content
							}
							return ""
						})

						// Stacks
						reStack := regexp.MustCompile(`@papan_tumpukan\s*\(\"(.*?)\"\)`)
						finalRaw = reStack.ReplaceAllStringFunc(finalRaw, func(m string) string {
							name := reStack.FindStringSubmatch(m)[1]
							if contents, ok := stacks[name]; ok {
								return strings.Join(contents, "\n")
							}
							return ""
						})

						raw = finalRaw
					}
				} else {
					// If no warisan, still try to resolve papan_tumpukan in the current file
					reStack := regexp.MustCompile(`@papan_tumpukan\s*\(\"(.*?)\"\)`)
					raw = reStack.ReplaceAllStringFunc(raw, func(m string) string {
						name := reStack.FindStringSubmatch(m)[1]
						if contents, ok := stacks[name]; ok {
							return strings.Join(contents, "\n")
						}
						return ""
					})
				}

				// 4. Javanese-Blade Compiler logic
				token := "kopong"
				if t, ok := sessions["_token"]; ok {
					token = t.Inspect()
				}
				csrfInput := fmt.Sprintf("<input type='hidden' name='_token' value='%s'>", token)
				raw = strings.ReplaceAll(raw, "@csrf", csrfInput)

				reEscaped := regexp.MustCompile(`{{\s*(.*?)\s*}}`)
				reRaw := regexp.MustCompile(`{!!\s*(.*?)\s*!!}`)
				reYen := regexp.MustCompile(`@yen\s*\((.*?)\)`)
				raw = reYen.ReplaceAllString(raw, " menowo $1 { ")
				raw = strings.ReplaceAll(raw, "@punkyan_yen", " } ")

				reTrack := regexp.MustCompile(`@track_neng\s*\((.*?)\s+neng\s+(.*?)\)`)
				raw = reTrack.ReplaceAllString(raw, " $_items = $2; $_len = dowo($_items); $_i = 0; track $_i < $_len { $1 = $_items[$_i]; ")
				raw = strings.ReplaceAll(raw, "@punkyan_track", " $_i = $_i + 1; } ")

				lines := strings.Split(raw, "\n")
				var buffer bytes.Buffer
				buffer.WriteString("$_out = ''; ")

				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if strings.HasPrefix(trimmed, "menowo") || strings.HasPrefix(trimmed, "track") || trimmed == "}" || strings.HasSuffix(trimmed, "}") || trimmed == "" {
						buffer.WriteString(line + "\n")
					} else {
						l := strings.ReplaceAll(line, "\"", "\\\"")
						l = reEscaped.ReplaceAllString(l, "\" + html_escape($1) + \"")
						l = reRaw.ReplaceAllString(l, "\" + string($1) + \"")
						buffer.WriteString("$_out = $_out + \"" + l + "\\n\";\n")
					}
				}
				buffer.WriteString("balekno $_out")

				env := object.NewEnvironment()
				for k, v := range data {
					env.Set(k, v)
				}

				l := lexer.New(buffer.String())
				p := parser.New(l)
				prog := p.ParseProgram()
				if len(p.Errors()) > 0 {
					return newError("render_template compile error: %s", p.Errors()[0])
				}

				return Eval(prog, env)
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
		"generate_token": {
			Fn: func(args ...object.Object) object.Object {
				b := make([]byte, 16)
				rand.Read(b)
				return &object.String{Value: hex.EncodeToString(b)}
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

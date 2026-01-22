package evaluator

import (
	"bufio"
	"fmt"
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

var builtins = map[string]*object.Builtin{
	"ketok": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"dowo": &object.Builtin{
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
	"isi": &object.Builtin{
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
	"takon": &object.Builtin{
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
	"moco_file": &object.Builtin{
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
	"nulis_file": &object.Builtin{
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
	"layani_web": &object.Builtin{
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
			fmt.Printf("🐺 Wolf404 melayani neng http://localhost%s\n", addr)

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

			err := http.ListenAndServe(addr, nil)
			if err != nil {
				return newError("Server web njeblug: %s", err.Error())
			}

			return NULL
		},
	},
}

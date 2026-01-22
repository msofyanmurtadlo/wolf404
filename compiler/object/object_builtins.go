package object

// ... existing code ...

const (
	// ...
	BUILTIN_OBJ = "BUILTIN"
)

// Builtin Function type
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// Module Object (e.g. 'http') is just a hash/map of functions? 
// Or a special object? Let's treat it as a Hash for simplicity or specialized Struct.

package object

const (
	CLASS_OBJ        = "CLASS"
	INSTANCE_OBJ     = "INSTANCE"
	BOUND_METHOD_OBJ = "BOUND_METHOD"
)

type Class struct {
	Name    string
	Super   *Class
	Methods map[string]*Function
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string {
	return "mold " + c.Name
}

type Instance struct {
	Class  *Class
	Fields map[string]Object
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string {
	return "instance of " + i.Class.Name
}

type BoundMethod struct {
	Method   *Function
	Instance *Instance
}

func (bm *BoundMethod) Type() ObjectType { return BOUND_METHOD_OBJ }
func (bm *BoundMethod) Inspect() string {
	return "bound method " + bm.Method.Inspect()
}

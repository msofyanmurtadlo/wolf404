package evaluator

import (
	"wolf404/compiler/object"
)

// ... existing code ...

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Class:
		instance := &object.Instance{Class: fn, Fields: make(map[string]object.Object)}
		if init, ok := fn.Methods["init"]; ok {
			// Call init with instance as $this
			val := applyMethod(init, args, instance)
			if isError(val) {
				return val
			}
		}
		return instance
	case *object.BoundMethod:
		return applyMethod(fn.Method, args, fn.Instance)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func applyMethod(fn *object.Function, args []object.Object, instance *object.Instance) object.Object {
	extendedEnv := extendFunctionEnv(fn, args)
	extendedEnv.Set("this", instance)
	evaluated := Eval(fn.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		// fmt.Printf("DEBUG: Binding param '%s' to %s\n", param.Value, args[paramIdx].Inspect())
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if obj == nil {
		return NULL
	}
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// ... evalHashLiteral, evalIndexExpression codes ... (I will append or ensure they exist)

package evaluator

import (
	"wolf404/compiler/ast"
	"wolf404/compiler/object"
)

func evalClassStatement(node *ast.ClassStatement, env *object.Environment) object.Object {
	methods := make(map[string]*object.Function)

	for _, stmt := range node.Body.Statements {
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			if funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
				// It's a method definition
				// Eval creates the function object
				obj := Eval(funcLit, env)
				if fn, ok := obj.(*object.Function); ok {
					methods[funcLit.Name] = fn
				}
			}
		}
	}

	var superClass *object.Class
	if node.SuperClass != nil {
		superObj, ok := env.Get(node.SuperClass.Value)
		if !ok {
			return newError("superclass not found: %s", node.SuperClass.Value)
		}
		superClass, ok = superObj.(*object.Class)
		if !ok {
			return newError("superclass must be a class: %s", superObj.Type())
		}
	}

	class := &object.Class{Name: node.Name.Value, Methods: methods, Super: superClass}
	env.Set(node.Name.Value, class)
	return NULL
}

func findMethod(class *object.Class, name string) *object.Function {
	if m, ok := class.Methods[name]; ok {
		return m
	}
	if class.Super != nil {
		return findMethod(class.Super, name)
	}
	return nil
}

func evalDotExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	rightIdent, ok := node.Right.(*ast.Identifier)
	if !ok {
		return newError("expected identifier as property name")
	}
	key := rightIdent.Value

	if instance, ok := left.(*object.Instance); ok {
		if val, ok := instance.Fields[key]; ok {
			return val
		}
		// Recursive method lookup
		if method := findMethod(instance.Class, key); method != nil {
			return &object.BoundMethod{Method: method, Instance: instance}
		}
	}

	return newError("property %s not found on %s", key, left.Type())
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	val := Eval(node.Right, env)
	if isError(val) {
		return val
	}

	if ident, ok := node.Left.(*ast.Identifier); ok {
		env.Set(ident.Value, val)
		return val
	}

	// Dot Assignment support ($obj.prop = val)
	if infixLeft, ok := node.Left.(*ast.InfixExpression); ok && infixLeft.Operator == "." {
		target := Eval(infixLeft.Left, env)
		if isError(target) {
			return target
		}

		propNameIdent, ok := infixLeft.Right.(*ast.Identifier)
		if !ok {
			return newError("expected property name")
		}

		if instance, ok := target.(*object.Instance); ok {
			instance.Fields[propNameIdent.Value] = val
			return val
		}
		return newError("cannot assign property to non-instance: %s", target.Type())
	}

	return newError("invalid assignment target")
}

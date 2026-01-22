package evaluator

import (
	"os"
	"wolf404/compiler/ast"
	"wolf404/compiler/lexer"
	"wolf404/compiler/object"
	"wolf404/compiler/parser"
)

func evalClassStatement(node *ast.ClassStatement, env *object.Environment) object.Object {
	methods := make(map[string]*object.Function)

	for _, stmt := range node.Body.Statements {
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			if funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
				// It's a method definition
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
		return newError("property access must be an identifier")
	}

	if instance, ok := left.(*object.Instance); ok {
		// Field access
		if val, ok := instance.Fields[rightIdent.Value]; ok {
			return val
		}

		// Method lookup
		if method := findMethod(instance.Class, rightIdent.Value); method != nil {
			return &object.BoundMethod{Instance: instance, Method: method}
		}

		return newError("property %s not found on INSTANCE", rightIdent.Value)
	}

	return newError("cannot access property of non-instance: %s", left.Type())
}

func evalAssignmentExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
	val := Eval(node.Right, env)
	if isError(val) {
		return val
	}

	// Case 1: Simple identifier variable assignment ($a = 1)
	if leftIdent, ok := node.Left.(*ast.Identifier); ok {
		env.Set(leftIdent.Value, val)
		return val
	}

	// Case 2: Property assignment (obj.prop = 1)
	if leftInfix, ok := node.Left.(*ast.InfixExpression); ok && leftInfix.Operator == "." {
		target := Eval(leftInfix.Left, env)
		if isError(target) {
			return target
		}

		propNameIdent, ok := leftInfix.Right.(*ast.Identifier)
		if !ok {
			return newError("property name must be identifier")
		}

		if instance, ok := target.(*object.Instance); ok {
			instance.Fields[propNameIdent.Value] = val
			return val
		}
		return newError("cannot assign property to non-instance: %s", target.Type())
	}

	return newError("invalid assignment target")
}

func evalSummonStatement(node *ast.SummonStatement, env *object.Environment) object.Object {
	path := node.Path.Value
	content, err := os.ReadFile(path)
	if err != nil {
		return newError("gagal ngundang file: %s", err.Error())
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		return newError("gagal ngonversi file nggih diundang: %v", p.Errors())
	}

	return Eval(program, env)
}

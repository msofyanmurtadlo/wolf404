package evaluator

import (
	"wolf404/compiler/ast"
	"wolf404/compiler/object"
)

func evalTrackStatement(ts *ast.TrackStatement, env *object.Environment) object.Object {
	var result object.Object

	const MAX_ITERATIONS = 1000000
	iterations := 0

	for {
		iterations++
		if iterations > MAX_ITERATIONS {
			return newError("Baleni kokean (Infinite Loop)! Watese 1 yuto iterasi.")
		}

		condition := Eval(ts.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ts.Body, env)

		// Handle return statements inside loop or errors
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

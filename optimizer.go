package main

/*
More optimzation techniques which can be added in future:
1. Common Subexpression Elimination
2. Loop Unrolling
3. Tail Call Optimization
4. Peephole Optimizations
*/

/* NOTE: remember to balance the trade-off between compilation time and runtime performance when adding more optimization passes */

import (
	"fmt"
	"strconv"
)

func optimize(ast ast) ast {
	optimized := constantFolding(ast)
	optimized = deadCodeElimination(optimized)
	optimized = functionInlining(optimized)
	return optimized
}

func constantFolding(a ast) ast {
	traverser(a, visitor{
		"CallExpression": func(n *node, p node) {
			if isArithmeticOperator(n.name) {
				if allParamsAreNumberLiterals(n.params) {
					result := evaluateArithmeticExpression(n)
					*n = node{
						kind:  "NumberLiteral",
						value: fmt.Sprintf("%v", result),
					}
				}
			}
		},
	})
	return a
}

func deadCodeElimination(a ast) ast {
	newBody := []node{}
	for _, n := range a.body {
		if !isDeadCode(n) {
			newBody = append(newBody, n)
		}
	}
	a.body = newBody
	return a
}

func functionInlining(a ast) ast {
		functions := make(map[string]node)

	traverser(a, visitor{
		"CallExpression": func(n *node, p node) {
			if n.name == "define" && len(n.params) == 3 {
				funcName := n.params[0].value
				funcBody := n.params[2]
				functions[funcName] = funcBody
			}
		},
	})

	traverser(a, visitor{
		"CallExpression": func(n *node, p node) {
			if funcBody, ok := functions[n.name]; ok {
				*n = funcBody
			}
		},
	})

	return a
}

func isArithmeticOperator(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/"
}

func allParamsAreNumberLiterals(params []node) bool {
	for _, param := range params {
		if param.kind != "NumberLiteral" {
			return false
		}
	}
	return true
}

func evaluateArithmeticExpression(n *node) float64 {
	op := n.name
	left, _ := strconv.ParseFloat(n.params[0].value, 64)
	right, _ := strconv.ParseFloat(n.params[1].value, 64)

	switch op {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	default:
		return 0
	}
}

func isDeadCode(n node) bool {
	return n.kind == "ExpressionStatement" && n.expression.kind == "NumberLiteral"
}

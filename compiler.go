package main

import (
	"fmt"
	"log"
	"strings"
)

type token struct {
	kind  string
	value string
}

func tokenizer(input string) ([]token, error) {
	input += "\n"
	current := 0
	tokens := []token{}

	for current < len([]rune(input)) {
		char := string([]rune(input)[current])

		if char == "(" {
			tokens = append(tokens, token{kind: "paren", value: "("})
			current++
			continue
		}

		if char == ")" {
			tokens = append(tokens, token{kind: "paren", value: ")"})
			current++
			continue
		}

		if char == " " {
			current++
			continue
		}

		if isNumber(char) {
			value := ""
			for isNumber(char) {
				value += char
				current++
				char = string([]rune(input)[current])
			}
			tokens = append(tokens, token{kind: "number", value: value})
			continue
		}

		if isLetter(char) {
			value := ""
			for isLetter(char) {
				value += char
				current++
				char = string([]rune(input)[current])
			}

			if value == "true" || value == "false" {
				tokens = append(tokens, token{kind: "boolean", value: value})
			} else {
				tokens = append(tokens, token{kind: "name", value: value})
			}
			continue
		}

		if char == `"` {
			value := ""
			current++
			char = string([]rune(input)[current])

			for char != `"` {
				value += char
				current++
				char = string([]rune(input)[current])
			}
			tokens = append(tokens, token{kind: "string", value: value})
			current++
			continue
		}

		if char == "+" || char == "-" || char == "*" || char == "/" {
			tokens = append(tokens, token{kind: "operator", value: char})
			current++
			continue
		}

		return nil, fmt.Errorf("unexpected character: %s at position %d", char, current)
	}

	return tokens, nil
}

func isNumber(char string) bool {
	if char == "" {
		return false
	}
	n := []rune(char)[0]
	return n >= '0' && n <= '9'
}

func isLetter(char string) bool {
	if char == "" {
		return false
	}
	n := []rune(char)[0]
	return (n >= 'a' && n <= 'z') || (n >= 'A' && n <= 'Z')
}

type node struct {
	kind       string
	value      string
	name       string
	operator   string
	left       *node
	right      *node
	callee     *node
	expression *node
	body       []node
	params     []node
	arguments  *[]node
	context    *[]node
}

type ast node

var pc int
var pt []token

func parser(tokens []token) (ast, error) {
	pc = 0
	pt = tokens
	ast := ast{kind: "Program", body: []node{}}

	for pc < len(pt) {
		n, err := walk()
		if err != nil {
			return ast, err
		}
		ast.body = append(ast.body, n)
	}

	return ast, nil
}

func walk() (node, error) {
	token := pt[pc]

	if token.kind == "number" {
		pc++
		return node{kind: "NumberLiteral", value: token.value}, nil
	}

	if token.kind == "string" {
		pc++
		return node{kind: "StringLiteral", value: token.value}, nil
	}

	if token.kind == "boolean" {
		pc++
		return node{kind: "BooleanLiteral", value: token.value}, nil
	}

	if token.kind == "operator" {
		pc++
		operator := token.value

		left, err := walk()
		if err != nil {
			return node{}, err
		}

		right, err := walk()
		if err != nil {
			return node{}, err
		}

		return node{
			kind:     "BinaryExpression",
			operator: operator,
			left:     &left,
			right:    &right,
		}, nil
	}

	if token.kind == "paren" && token.value == "(" {
		pc++
		token = pt[pc]

		n := node{kind: "CallExpression", name: token.value, params: []node{}}
		pc++
		token = pt[pc]

		for token.kind != "paren" || (token.kind == "paren" && token.value != ")") {
			param, err := walk()
			if err != nil {
				return node{}, err
			}
			n.params = append(n.params, param)
			token = pt[pc]
		}

		pc++
		return n, nil
	}

	return node{}, fmt.Errorf("unexpected token: %s", token.kind)
}

type visitor map[string]func(n *node, p node)

func traverser(a ast, v visitor) {
	traverseNode(node(a), node{}, v)
}

func traverseArray(a []node, p node, v visitor) {
	for _, child := range a {
		traverseNode(child, p, v)
	}
}

func traverseNode(n, p node, v visitor) {
	if va, ok := v[n.kind]; ok {
		va(&n, p)
	}

	switch n.kind {
	case "Program":
		traverseArray(n.body, n, v)
	case "CallExpression":
		traverseArray(n.params, n, v)
	case "BinaryExpression":
		traverseNode(*n.left, n, v)
		traverseNode(*n.right, n, v)
	}
}

func transformer(a ast) ast {
	nast := ast{kind: "Program", body: []node{}}
	a.context = &nast.body

	traverser(a, map[string]func(n *node, p node){
		"NumberLiteral": func(n *node, p node) {
			*p.context = append(*p.context, node{kind: "NumberLiteral", value: n.value})
		},
		"StringLiteral": func(n *node, p node) {
			*p.context = append(*p.context, node{kind: "StringLiteral", value: n.value})
		},
		"BooleanLiteral": func(n *node, p node) {
			*p.context = append(*p.context, node{kind: "BooleanLiteral", value: n.value})
		},
		"BinaryExpression": func(n *node, p node) {
			left := node{
				kind: "BinaryExpression",
				left: n.left,
			}
			right := node{
				kind: "BinaryExpression",
				right: n.right,
			}
			*p.context = append(*p.context, node{
				kind:     "BinaryExpression",
				operator: n.operator,
				left:     &left,
				right:    &right,
			})
		},
		"CallExpression": func(n *node, p node) {
			e := node{
				kind: "CallExpression",
				callee: &node{
					kind: "Identifier",
					name: n.name,
				},
				arguments: new([]node),
			}
			n.context = e.arguments
			if p.kind != "CallExpression" {
				es := node{
					kind:       "ExpressionStatement",
					expression: &e,
				}
				*p.context = append(*p.context, es)
			} else {
				*p.context = append(*p.context, e)
			}
		},
	})

	return nast
}

func codeGenerator(n node) string {
	switch n.kind {
	case "Program":
		var r []string
		for _, no := range n.body {
			r = append(r, codeGenerator(no))
		}
		return strings.Join(r, "\n")
	case "ExpressionStatement":
		return codeGenerator(*n.expression) + ";"
	case "CallExpression":
		var ra []string
		c := codeGenerator(*n.callee)
		for _, no := range *n.arguments {
			ra = append(ra, codeGenerator(no))
		}
		return c + "(" + strings.Join(ra, ", ") + ")"
	case "Identifier":
		return n.name
	case "NumberLiteral":
		return n.value
	case "StringLiteral":
		return `"` + n.value + `"`
	case "BooleanLiteral":
		return n.value
	case "BinaryExpression":
		return codeGenerator(*n.left) + " " + n.operator + " " + codeGenerator(*n.right)
	default:
		log.Fatal("unknown node kind")
		return ""
	}
}

func compiler(input string) (string, error) {
	tokens, err := tokenizer(input)
	if err != nil {
		return "", err
	}
	ast, err := parser(tokens)
	if err != nil {
		return "", err
	}
	optimizedAst := optimize(ast)
	nast := transformer(optimizedAst)
	out := codeGenerator(node(nast))
	return out, nil
}

func main() {
	program := `(+ 5 (- 3 2))`
	out, err := compiler(program)
	if err != nil {
		log.Fatalf("Compilation error: %v", err)
	}
	fmt.Println(out)
}

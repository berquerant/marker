package main

type (
	Expr interface {
		IsExpr()
	}
	Operator interface {
		IsOperator()
	}
	AddExpr struct{}
)

func main() {
	x := &AddExpr{}
	var expr Expr
	expr = x
	println(expr)
	var op Operator
	op = x
	println(op)
}

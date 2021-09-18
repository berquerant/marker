package main

type (
	Expr interface {
		IsExpr()
	}
	StatementExpr struct{}
)

func main() {
	var expr Expr
	expr = &StatementExpr{}
	println(expr)
}

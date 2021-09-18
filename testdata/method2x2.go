package main

type (
	Expr interface {
		IsExpr()
	}
	Node interface {
		IsNode()
	}
	FlexExpr  struct{}
	FixedExpr struct{}
)

func main() {
	a := &FlexExpr{}
	b := &FixedExpr{}
	var (
		expr Expr
		node Node
	)
	expr = a
	expr = b
	node = a
	node = b
	println(expr, node)
}

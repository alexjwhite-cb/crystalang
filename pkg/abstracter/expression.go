package abstracter

type ParenExpression struct {
	Lparen Position
	X      Expr
	Rparen Position
}

func (p *ParenExpression) Pos() int { return p.Lparen.Start }
func (p *ParenExpression) End() int { return p.Rparen.Start }

func (p *ParenExpression) exprNode() {}

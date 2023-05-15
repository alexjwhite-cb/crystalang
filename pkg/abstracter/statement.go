package abstracter

type BlockStatement struct {
	Lbrace *Position
	List   []Stmt
	Rbrace *Position
}

func (b *BlockStatement) Pos() int { return b.Lbrace.Start }
func (b *BlockStatement) End() int { return b.Rbrace.Start }

func (*BlockStatement) stmtNode() {}

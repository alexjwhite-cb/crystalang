package abstracter

type BlockStatement struct {
	Lbrace   *Position
	Identity string
	List     []Stmt
	RBrace   *Position
}

func (b *BlockStatement) Pos() int { return b.Lbrace.Start }
func (b *BlockStatement) End() int { return b.Lbrace.Start }

func (*BlockStatement) stmtNode() {}

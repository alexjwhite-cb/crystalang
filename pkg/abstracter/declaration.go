package abstracter

type MethodDeclaration struct {
	Meth *Ident
	Name string
	Loc  Position
}

func (m *MethodDeclaration) Pos() int { return m.Loc.Start }
func (m *MethodDeclaration) End() int { return m.Loc.End }

type DescriptorDeclaration struct {
	Describe *Ident
	Name     string
	Loc      Position
}

func (d *DescriptorDeclaration) Pos() int { return d.Loc.Start }
func (d *DescriptorDeclaration) End() int { return d.Loc.End }

type ObjectDeclaration struct {
	Object *Ident
	Name   string
	Loc    Position
}

func (o *ObjectDeclaration) Pos() int { return o.Loc.Start }
func (o *ObjectDeclaration) End() int { return o.Loc.End }

type BadDeclaration struct {
	From, To int
}

func (b *BadDeclaration) Pos() int { return b.From }
func (b *BadDeclaration) End() int { return b.To }

func (*BadDeclaration) declNode()        {}
func (*MethodDeclaration) declNode()     {}
func (*DescriptorDeclaration) declNode() {}
func (*ObjectDeclaration) declNode()     {}

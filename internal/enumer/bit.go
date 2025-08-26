package enumer

func (g *Generator) makeBitwize() {
	if !g.flags.bitwise {
		return
	}

	g.data.Bitwise = true
	//todo: int values check
}

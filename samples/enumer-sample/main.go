package main

import fsenum "enumer_sample/formatstyleenum"

func main() {
	fs := fsenum.Bold | fsenum.Italic
	if fs.Has(fsenum.Bold) {
		println("Bold is set")
	}
}

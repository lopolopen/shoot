package shoot

type TmplDataBase struct {
	Cmd         string
	PackageName string
	TypeName    string
}

func NewTmplDataBase() *TmplDataBase {
	d := &TmplDataBase{}
	return d
}

func (d *TmplDataBase) SetCmd(cmd string) {
	d.Cmd = cmd
}

func (d *TmplDataBase) SetTypeName(typeName string) {
	d.TypeName = typeName
}

func (d *TmplDataBase) SetPackageName(pkgName string) {
	d.PackageName = pkgName
}

// func (d *TmplDataBase) String() string {
// 	return fmt.Sprintf(`Cmd: %s, PackageName: %s, TypeName: %s`, d.Cmd, d.PackageName, d.TypeName)
// }

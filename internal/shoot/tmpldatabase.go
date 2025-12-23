package shoot

type TmplDataBase struct {
	CmdLine     string
	Version     string
	PackageName string
	TypeName    string
}

func NewTmplDataBase(cmdline, version string) *TmplDataBase {
	d := &TmplDataBase{
		CmdLine: cmdline,
		Version: version,
	}
	return d
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

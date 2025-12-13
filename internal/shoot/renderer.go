package shoot

// func Render(dir, outputName, subCmd, tmplTxt string, funcs template.FuncMap, data interface{}) {
// 	var buff bytes.Buffer

// 	tmpl, err := template.New(subCmd).Funcs(funcs).Parse(tmplTxt)
// 	if err != nil {
// 		log.Fatalf("❌ parsing template: %s", err)
// 	}
// 	err = tmpl.Execute(&buff, data)
// 	if err != nil {
// 		log.Fatalf("❌ executing template: %s", err)
// 	}

// 	notedown(dir, outputName, buff.Bytes())
// }

// func notedown(dir, outputName string, src []byte) {
// 	// Write to tmpfile first
// 	tmpFile, err := os.CreateTemp(dir, fmt.Sprintf(".%s_", outputName))
// 	if err != nil {
// 		log.Fatalf("❌ creating temporary file for output: %s", err)
// 	}
// 	_, err = tmpFile.Write(src)
// 	if err != nil {
// 		tmpFile.Close()
// 		os.Remove(tmpFile.Name())
// 		log.Fatalf("❌ writing output: %s", err)
// 	}
// 	tmpFile.Close()

// 	// Rename tmpfile to output file
// 	err = os.Rename(tmpFile.Name(), outputName)
// 	if err != nil {
// 		log.Fatalf("❌ moving tempfile to output file: %s", err)
// 	}
// }

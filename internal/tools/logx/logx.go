package logx

import (
	"encoding/json"
	"log"
)

func Fatal(v ...any) {
	xs := append([]any{"❌"}, v...)
	log.Fatal(xs...)
}

func Fatalf(format string, v ...any) {
	log.Fatalf("❌ "+format, v...)
}

func Warn(v ...any) {
	xs := append([]any{"⚠️"}, v...)
	log.Print(xs...)
}

func Warnf(format string, v ...any) {
	log.Printf("⚠️ "+format, v...)
}

func Pinln(v ...any) {
	xs := append([]any{"📌"}, v...)
	log.Println(xs...)
}

func Debugln(v ...any) {
	xs := append([]any{"🐛"}, v...)
	log.Println(xs...)
}

func DebugJSONln(v ...any) {
	xs := []any{"🐛"}
	for _, x := range v {
		str, ok := x.(string)
		if ok {
			xs = append(xs, str)
			continue
		}
		j, _ := json.MarshalIndent(x, "", "  ")
		xs = append(xs, string(j))
	}
	log.Println(xs...)
}

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

func Warn(v ...any) { //width: A not W
	xs := append([]any{"⚠️", " "}, v...)
	log.Print(xs...)
}

func Warnf(format string, v ...any) {
	log.Printf("⚠️ "+" "+format, v...)
}

func Pin(v ...any) {
	xs := append([]any{"📌"}, v...)
	log.Println(xs...)
}

func PinJSON(v ...any) {
	xs := []any{"📌"}
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

func Debug(v ...any) {
	xs := append([]any{"🐛"}, v...)
	log.Println(xs...)
}

func DebugJSON(v ...any) {
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

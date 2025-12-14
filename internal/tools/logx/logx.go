package logx

import (
	"encoding/json"
	"fmt"
	"log"
)

func Fatal(v ...any) {
	s := fmt.Sprint(v...)
	log.Fatalf("❌ %s", s)
}

func Fatalf(format string, v ...any) {
	log.Fatalf("❌ "+format, v...)
}

func Warn(v ...any) {
	s := fmt.Sprint(v...)
	log.Printf("⚠️ %s", s)
}

func Warnf(format string, v ...any) {
	log.Printf("⚠️ "+format, v...)
}

func Pinln(v ...any) {
	s := fmt.Sprint(v...)
	log.Printf("📌 %s\n", s)
}

func Debugln(v ...any) {
	s := fmt.Sprint(v...)
	log.Printf("🐛 %s\n", s)
}

func DebugJSONln(v ...any) {
	var xs []any
	for _, x := range v {
		str, ok := x.(string)
		if ok {
			xs = append(xs, str)
			continue
		}
		j, _ := json.MarshalIndent(x, "", "  ")
		xs = append(xs, string(j))
	}
	s := fmt.Sprint(xs...)
	log.Printf("🐛 %s\n", s)
}

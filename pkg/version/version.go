package version

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

func CreateHealthHandler(service string) http.Handler {
	values := readBuildVersion(service)
	rawJSON, _ := json.MarshalIndent(values, "", "  ")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(rawJSON)
	})
}

func readBuildVersion(service string) map[string]string {
	m := make(map[string]string)
	m["service"] = service

	file, err := os.Open("./build.version")
	if err != nil {
		m["error"] = err.Error()
		return m
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		index := strings.IndexByte(line, '=')
		if index <= 0 {
			continue
		}
		key := strings.TrimSpace(line[0:index])
		value := strings.TrimSpace(line[index+1:])
		if key == "" {
			continue
		}
		m[key] = value
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return m
}

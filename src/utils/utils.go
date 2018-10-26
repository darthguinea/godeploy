package utils

import (
	"os"

	"../log"
)

func LoadTemplate(path string) string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 1048576)
	x, _ := file.Read(data)

	return string(data[:x])
}

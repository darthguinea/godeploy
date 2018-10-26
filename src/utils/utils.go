package utils

import (
	"os"

	"../log"
)

func LoadTemplate(path string) string {
	FILE_SIZE := 1048576
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, FILE_SIZE)
	x, _ := file.Read(data)

	return string(data[:x])
}

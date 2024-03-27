package helpers

import (
	"encoding/json"
	"os"
)

type ConnFile struct {
	ConnectionString string
}

func GetConnString(connFilePath string) (string, error) {
	buf, err := os.ReadFile(connFilePath)

	connFile := ConnFile{}

	if err != nil {
		return "", err
	}

	err = json.Unmarshal(buf, &connFile)

	return connFile.ConnectionString, err
}
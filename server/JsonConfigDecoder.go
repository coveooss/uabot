package server

import (
	"encoding/json"
	"io"
	"github.com/adambbolduc/uabot/explorerlib"
)

func DecodeConfig(jsonReader io.Reader) (*explorerlib.Config, error){
	config := &explorerlib.Config{}
	err := json.NewDecoder(jsonReader).Decode(config)
	return config, err
}


package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StatusResponse struct {
	Online  bool `json:"online"`
	Players struct {
		Online uint64 `json:"online"`
		Max    uint64 `json:"max"`
	} `json:"players"`
	Icon *string `json:"icon"`
}

func GetServerStatus(server *UniqueServer) (*StatusResponse, error) {
	host := server.Hostname

	if (server.Type == "java" && server.Port != 25565) || (server.Type == "bedrock" && server.Port != 19132) {
		host += fmt.Sprintf(":%d", server.Port)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status/%s/%s", config.PingServerHost, server.Type, host), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", config.AuthToken)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ping: unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var result StatusResponse

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

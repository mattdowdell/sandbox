package dashboards

import (
	"encoding/json"

	"github.com/mattdowdell/sandbox/pkg/dashboards/connectrpc"
)

func ConnectRPC() (string, error) {
	board, err := connectrpc.NewDashboard()
	if err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(board, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

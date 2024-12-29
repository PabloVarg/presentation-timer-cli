package api

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type ErrorResponse struct {
	Error    string              `json:"error"`
	Messages map[string][]string `json:"messages"`
}

func ExtractErrorMsg(body io.ReadCloser) error {
	var res ErrorResponse

	if err := json.NewDecoder(body).Decode(&res); err != nil {
		return err
	}

	if res.Messages != nil {
		return fmt.Errorf(
			"%s",
			formatMessages(res.Messages),
		)
	}

	return fmt.Errorf("%s", res.Error)
}

func formatMessages(messages map[string][]string) string {
	labels := make([]string, 0, 1)

	for _, labelMessages := range messages {
		labels = append(labels,
			strings.Join(
				labelMessages,
				"\n",
			),
		)
	}

	return strings.Join(labels, "\n")
}

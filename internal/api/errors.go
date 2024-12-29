package api

import (
	"encoding/json"
	"fmt"
	"io"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func ExtractErrorMsg(body io.ReadCloser) error {
	var res ErrorResponse

	if err := json.NewDecoder(body).Decode(&res); err != nil {
		return err
	}

	return fmt.Errorf("%s", res.Error)
}

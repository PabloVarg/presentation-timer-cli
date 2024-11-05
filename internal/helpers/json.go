package helpers

import (
	"encoding/json"
	"io"
)

func ReadJSON(r io.Reader, output any) error {
	if err := json.NewDecoder(r).Decode(output); err != nil {
		return err
	}

	return nil
}

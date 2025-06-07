package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func readRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		fmt.Errorf("json.Decode: %w", err)
		return err
	}

	return nil
}

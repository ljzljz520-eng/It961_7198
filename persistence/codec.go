package persistence

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func encode(value any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		return nil, fmt.Errorf("encode record: %w", err)
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return fmt.Errorf("empty record")
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode record: %w", err)
	}
	return nil
}

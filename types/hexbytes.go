package types

import (
	"encoding/hex"
	"encoding/json"
)

type HexBytes []byte

func (h *HexBytes) UnmarshalJSON(b []byte) error {
	// b viene como: `"4a6f686e"`
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	decoded, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	*h = decoded
	return nil
}

package provider

import (
	"encoding/json"
	"io"
)

// decodeJSONBody is a tiny helper used by mock-server tests to inspect request
// payloads. Kept in a _test.go file so it never leaks into production binaries.
func decodeJSONBody(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

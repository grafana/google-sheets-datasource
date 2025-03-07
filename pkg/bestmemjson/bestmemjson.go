// Package bestmemjson provides JSON operations using the most memory-efficient implementation
// based on benchmark testing. It uses encoding/json for Marshal operations and
// jsoniter for Unmarshal operations.
package bestmemjson

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
)

// Marshal encodes the input using encoding/json Marshal which uses less memory
// than jsoniter for Marshal operations (112B vs 120B per op)
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes the input using jsoniter Unmarshal which uses less memory
// than encoding/json for Unmarshal operations (200B vs 352B per op)
func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

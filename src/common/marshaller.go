package common

import (
	"encoding/json"
)

// Marshaller encodes and decodes objects to and from byte arrays
type Marshaller struct{}

// Marshall ..
func (m *Marshaller) Marshall(data interface{}) []byte {
	pkt, _ := json.Marshal(data)
	return pkt
}

// Unmarshall ...
func (m *Marshaller) Unmarshall(data []byte, result interface{}) error {

	return json.Unmarshal(data, result)
}

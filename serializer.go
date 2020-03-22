package flamingo

import (
	"encoding/json"
	"fmt"
)

// DefaultSerializer is used when no serializer is applied to the
// current app, route or response.
var DefaultSerializer *Serializer = JSONSerializer()

// Serializer is used for serializing response data and de-serializing
// input data.
type Serializer struct {
	// Determine if this is a raw single key serializer
	IsRaw bool
	// The key for this serializer, if it is a Raw serializer.
	Key string
	// The content type to use when serializing response data.
	ContentType string
	// Serialize should take a map as input and serialize it to a
	// string. If any issues are encountered, return an empty
	// string and an error.
	Serialize func(map[string]interface{}) (string, error)
	// Serialize should take a string as input and serialize it to a
	// map. If any issues are encountered, return a nil map and an
	// error.
	Deserialize func(string) (map[string]interface{}, error)
}

// MakeRawData returns the Raw Data for the specified data using the
// key configured in this Serializer.
func (serializer *Serializer) MakeRawData(data string) map[string]interface{} {
	res := map[string]interface{}{}
	if serializer.IsRaw {
		res[serializer.Key] = data
	}

	return res
}

// JSONSerializer returns a Serializer for JSON serialization.
func JSONSerializer() *Serializer {
	return &Serializer{
		ContentType: "application/json",
		Serialize: func(data map[string]interface{}) (string, error) {
			encoded, err := json.Marshal(data)
			if err != nil {
				return "", err
			}

			return string(encoded), nil
		},
		Deserialize: func(data string) (map[string]interface{}, error) {
			res := map[string]interface{}{}
			err := json.Unmarshal([]byte(data), &res)
			if err != nil {
				return nil, err
			}

			return res, nil
		},
	}
}

// DownloadSerializer returns a Serializer for file downloads.
func DownloadSerializer() *Serializer {
	return NewRawSerializer("data", "application/octet-stream")
}

// TextSerializer returns a Serializer for Plain Text.
func TextSerializer() *Serializer {
	return NewRawSerializer("text", "text/plain")
}

// NewRawSerializer creates a serializer that takes raw input and
// applies the specified content type to responses that are serialized
// using it. When serializing data with this serializer, you should
// use a map with one entry being the data to serialize mapped to the
// key specified in key.
func NewRawSerializer(key string, contentType string) *Serializer {
	return &Serializer{
		IsRaw:       true,
		Key:         key,
		ContentType: contentType,
		Serialize: func(data map[string]interface{}) (string, error) {
			if out, exists := data[key]; exists {
				return out.(string), nil
			}

			return "", fmt.Errorf("invalid %s data", key)
		},
		Deserialize: func(data string) (map[string]interface{}, error) {
			return map[string]interface{}{
				key: data,
			}, nil
		},
	}
}

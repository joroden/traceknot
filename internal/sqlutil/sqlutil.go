package sqlutil

import "encoding/json"

func Placeholders(count int) string {
	output := make([]byte, 0, count*2)
	for index := 0; index < count; index++ {
		if index > 0 {
			output = append(output, ',')
		}
		output = append(output, '?')
	}
	return string(output)
}

func ToAnySlice(values []string) []any {
	output := make([]any, len(values))
	for index, value := range values {
		output[index] = value
	}
	return output
}

func MustJSON(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

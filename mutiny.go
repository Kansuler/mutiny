package mutiny

import (
	"encoding/json"
	"reflect"
	"regexp"
	"strings"
)

// pickValuesFromPossibleValues picks all the possible values from a PossibleValues struct.
func pickValuesFromPossibleValues(originalValues []any, possibleValues any) []any {
	// Process the field if it can be typecast to PossibleValues.
	if possibleValues, ok := possibleValues.(PossibleValues); ok {
		// Fill the flatMapWithdrawalValues with the selection for this possible values.
		switch possibleValues.use {
		case Fail:
			return append(originalValues, possibleValues.Fail...)
		case Erroneous:
			return append(originalValues, possibleValues.Erroneous...)
		case Nil:
			return append(originalValues, nil)
		default:
			// Use Pass values by default
			if len(possibleValues.Pass) > 0 {
				return append(originalValues, possibleValues.Pass...)
			}
		}
	}
	return originalValues
}

// createFlatMap creates a flat map of all the possible values for each field in a payload.
func createFlatMap(value reflect.Value) map[string][]any {
	flatMap := make(map[string][]any, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		// If field doesn't exist, initiate it
		if _, ok := flatMap[value.Type().Field(i).Name]; !ok {
			flatMap[value.Type().Field(i).Name] = make([]any, 0)
		}

		flatMap[value.Type().Field(i).Name] = pickValuesFromPossibleValues(
			flatMap[value.Type().Field(i).Name],
			value.Field(i).Interface(),
		)
	}

	return flatMap
}

// buildByField builds a slice of maps by field.
func buildByField(inputMaps []map[string]any, field string, values []any) []map[string]any {
	if len(values) == 0 {
		return inputMaps
	}

	var resultMaps []map[string]any
	for _, value := range values {
		if len(inputMaps) == 0 {
			var result = make(map[string]any)
			result[field] = value
			resultMaps = append(resultMaps, result)
			continue
		}

		for _, originalMap := range inputMaps {
			result := make(map[string]any)
			for k, v := range originalMap {
				result[k] = v
			}

			result[field] = value
			resultMaps = append(resultMaps, result)
		}
	}

	return resultMaps
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts a string from camelCase to snake_case.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// TestUnit is a unit of fields marshalled from the payload together with a JSON byte representation for a request body.
type TestUnit struct {
	Fields      map[string]any
	RequestBody []byte
}

// config is a struct that contains the configuration to be executed during serialization of payload.
type config struct {
	FieldFormatter func(string) string
}

// Option is a signature for all the with* functions that modify the configuration.
type Option func(*config)

// WithFieldFormatter is a function that modifies the configuration to use a custom field formatter.
func WithFieldFormatter(formatter func(string) string) Option {
	return func(cfg *config) {
		cfg.FieldFormatter = formatter
	}
}

// Riot is a function that serializes a payload into a slice of TestUnits.
func Riot[T any](payload T, opts ...Option) ([]TestUnit, error) {
	// Default configuration
	config := &config{
		FieldFormatter: ToSnakeCase,
	}

	// Loop through and apply any options sent in as argument.
	for _, opt := range opts {
		opt(config)
	}

	// Create a flat map of all the possible values for each field in a withdrawal and bank_account.
	flatMapPayloadValues := createFlatMap(reflect.ValueOf(payload))

	// Build the withdrawal payload by field.
	var payloadMaps []map[string]any
	for field, values := range flatMapPayloadValues {
		payloadMaps = buildByField(payloadMaps, config.FieldFormatter(field), values)
	}

	// JSON serialize all the withdrawal payloads.
	var resultPayloads []TestUnit
	for _, result := range payloadMaps {
		requestBody, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		resultPayloads = append(resultPayloads, TestUnit{Fields: result, RequestBody: requestBody})
	}

	return resultPayloads, nil
}

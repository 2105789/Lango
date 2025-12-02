package main

import (
	"encoding/json"
	"fmt"
)

// NativeJsonParse is a built-in function to parse JSON strings
type NativeJsonParse struct{}

func (n *NativeJsonParse) Arity() int {
	return 1
}

func (n *NativeJsonParse) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("jsonParse requires exactly 1 argument")
	}

	jsonStr, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("jsonParse: argument must be a string")
	}

	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, fmt.Errorf("jsonParse: invalid JSON - %v", err)
	}

	return convertGoToLango(result), nil
}

func (n *NativeJsonParse) String() string {
	return "<native fn jsonParse>"
}

// NativeJsonStringify is a built-in function to convert Lango values to JSON
type NativeJsonStringify struct{}

func (n *NativeJsonStringify) Arity() int {
	return 1
}

func (n *NativeJsonStringify) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("jsonStringify requires exactly 1 argument")
	}

	goValue := convertLangoToGo(arguments[0])

	jsonBytes, err := json.Marshal(goValue)
	if err != nil {
		return nil, fmt.Errorf("jsonStringify: failed to serialize - %v", err)
	}

	return string(jsonBytes), nil
}

func (n *NativeJsonStringify) String() string {
	return "<native fn jsonStringify>"
}

// NativeJsonStringifyPretty is a built-in function to convert Lango values to pretty-printed JSON
type NativeJsonStringifyPretty struct{}

func (n *NativeJsonStringifyPretty) Arity() int {
	return 1
}

func (n *NativeJsonStringifyPretty) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("jsonStringifyPretty requires exactly 1 argument")
	}

	goValue := convertLangoToGo(arguments[0])

	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("jsonStringifyPretty: failed to serialize - %v", err)
	}

	return string(jsonBytes), nil
}

func (n *NativeJsonStringifyPretty) String() string {
	return "<native fn jsonStringifyPretty>"
}

// convertGoToLango converts Go JSON types to Lango types
func convertGoToLango(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		langoMap := make(map[string]interface{})
		for key, val := range v {
			langoMap[key] = convertGoToLango(val)
		}
		return langoMap
	case []interface{}:
		langoArray := make([]interface{}, len(v))
		for i, val := range v {
			langoArray[i] = convertGoToLango(val)
		}
		return langoArray
	case float64, string, bool, nil:
		return v
	default:
		return v
	}
}

// convertLangoToGo converts Lango types to Go JSON-compatible types
func convertLangoToGo(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		goMap := make(map[string]interface{})
		for key, val := range v {
			goMap[key] = convertLangoToGo(val)
		}
		return goMap
	case []interface{}:
		goArray := make([]interface{}, len(v))
		for i, val := range v {
			goArray[i] = convertLangoToGo(val)
		}
		return goArray
	case *LangoInstance:
		// Convert instance fields to map
		goMap := make(map[string]interface{})
		for key, val := range v.fields {
			goMap[key] = convertLangoToGo(val)
		}
		return goMap
	case float64, string, bool, nil:
		return v
	default:
		// Attempt to convert to string as fallback
		return fmt.Sprintf("%v", v)
	}
}

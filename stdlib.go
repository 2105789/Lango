package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Helper functions for stdlib

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	if f, ok := object.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	if callable, ok := object.(LangoCallable); ok {
		return callable.String()
	}
	if instance, ok := object.(*LangoInstance); ok {
		return instance.String()
	}
	if array, ok := object.([]interface{}); ok {
		var builder strings.Builder
		builder.WriteString("[")
		for idx, val := range array {
			if idx > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(stringify(val))
		}
		builder.WriteString("]")
		return builder.String()
	}
	return fmt.Sprintf("%v", object)
}

// String Functions

type NativeStrLen struct{}

func (n *NativeStrLen) Arity() int { return 1 }

func (n *NativeStrLen) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	str, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("strLen: argument must be a string")
	}
	return float64(len(str)), nil
}

func (n *NativeStrLen) String() string { return "<native fn strLen>" }

type NativeStrUpper struct{}

func (n *NativeStrUpper) Arity() int { return 1 }

func (n *NativeStrUpper) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	str, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("strUpper: argument must be a string")
	}
	return strings.ToUpper(str), nil
}

func (n *NativeStrUpper) String() string { return "<native fn strUpper>" }

type NativeStrLower struct{}

func (n *NativeStrLower) Arity() int { return 1 }

func (n *NativeStrLower) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	str, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("strLower: argument must be a string")
	}
	return strings.ToLower(str), nil
}

func (n *NativeStrLower) String() string { return "<native fn strLower>" }

type NativeStrSplit struct{}

func (n *NativeStrSplit) Arity() int { return 2 }

func (n *NativeStrSplit) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	str, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("strSplit: first argument must be a string")
	}
	delimiter, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("strSplit: second argument must be a string")
	}
	parts := strings.Split(str, delimiter)
	result := make([]interface{}, len(parts))
	for i, part := range parts {
		result[i] = part
	}
	return result, nil
}

func (n *NativeStrSplit) String() string { return "<native fn strSplit>" }

type NativeStrJoin struct{}

func (n *NativeStrJoin) Arity() int { return 2 }

func (n *NativeStrJoin) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	arr, ok := arguments[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("strJoin: first argument must be an array")
	}
	delimiter, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("strJoin: second argument must be a string")
	}
	parts := make([]string, len(arr))
	for i, item := range arr {
		parts[i] = stringify(item)
	}
	return strings.Join(parts, delimiter), nil
}

func (n *NativeStrJoin) String() string { return "<native fn strJoin>" }

type NativeStrContains struct{}

func (n *NativeStrContains) Arity() int { return 2 }

func (n *NativeStrContains) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	str, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("strContains: first argument must be a string")
	}
	substr, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("strContains: second argument must be a string")
	}
	return strings.Contains(str, substr), nil
}

func (n *NativeStrContains) String() string { return "<native fn strContains>" }

type NativeStrReplace struct{}

func (n *NativeStrReplace) Arity() int { return 3 }

func (n *NativeStrReplace) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	str, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("strReplace: first argument must be a string")
	}
	old, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("strReplace: second argument must be a string")
	}
	new, ok := arguments[2].(string)
	if !ok {
		return nil, fmt.Errorf("strReplace: third argument must be a string")
	}
	return strings.ReplaceAll(str, old, new), nil
}

func (n *NativeStrReplace) String() string { return "<native fn strReplace>" }

type NativeStrTrim struct{}

func (n *NativeStrTrim) Arity() int { return 1 }

func (n *NativeStrTrim) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	str, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("strTrim: argument must be a string")
	}
	return strings.TrimSpace(str), nil
}

func (n *NativeStrTrim) String() string { return "<native fn strTrim>" }

// Array Functions

type NativeArrayLen struct{}

func (n *NativeArrayLen) Arity() int { return 1 }

func (n *NativeArrayLen) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	arr, ok := arguments[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("arrayLen: argument must be an array")
	}
	return float64(len(arr)), nil
}

func (n *NativeArrayLen) String() string { return "<native fn arrayLen>" }

type NativeArrayPush struct{}

func (n *NativeArrayPush) Arity() int { return 2 }

func (n *NativeArrayPush) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	arr, ok := arguments[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("arrayPush: first argument must be an array")
	}
	arr = append(arr, arguments[1])
	return arr, nil
}

func (n *NativeArrayPush) String() string { return "<native fn arrayPush>" }

type NativeArrayPop struct{}

func (n *NativeArrayPop) Arity() int { return 1 }

func (n *NativeArrayPop) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	arr, ok := arguments[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("arrayPop: argument must be an array")
	}
	if len(arr) == 0 {
		return nil, fmt.Errorf("arrayPop: cannot pop from empty array")
	}
	return arr[len(arr)-1], nil
}

func (n *NativeArrayPop) String() string { return "<native fn arrayPop>" }

type NativeArraySlice struct{}

func (n *NativeArraySlice) Arity() int { return 3 }

func (n *NativeArraySlice) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	arr, ok := arguments[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("arraySlice: first argument must be an array")
	}
	startFloat, ok := arguments[1].(float64)
	if !ok {
		return nil, fmt.Errorf("arraySlice: second argument must be a number")
	}
	endFloat, ok := arguments[2].(float64)
	if !ok {
		return nil, fmt.Errorf("arraySlice: third argument must be a number")
	}
	start := int(startFloat)
	end := int(endFloat)
	if start < 0 || start > len(arr) || end < start || end > len(arr) {
		return nil, fmt.Errorf("arraySlice: invalid slice bounds")
	}
	return arr[start:end], nil
}

func (n *NativeArraySlice) String() string { return "<native fn arraySlice>" }

type NativeArrayConcat struct{}

func (n *NativeArrayConcat) Arity() int { return 2 }

func (n *NativeArrayConcat) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	arr1, ok := arguments[0].([]interface{})
	if !ok {
		return nil, fmt.Errorf("arrayConcat: first argument must be an array")
	}
	arr2, ok := arguments[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("arrayConcat: second argument must be an array")
	}
	result := make([]interface{}, 0, len(arr1)+len(arr2))
	result = append(result, arr1...)
	result = append(result, arr2...)
	return result, nil
}

func (n *NativeArrayConcat) String() string { return "<native fn arrayConcat>" }

// Math Functions

type NativeMathFloor struct{}

func (n *NativeMathFloor) Arity() int { return 1 }

func (n *NativeMathFloor) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	num, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathFloor: argument must be a number")
	}
	return math.Floor(num), nil
}

func (n *NativeMathFloor) String() string { return "<native fn mathFloor>" }

type NativeMathCeil struct{}

func (n *NativeMathCeil) Arity() int { return 1 }

func (n *NativeMathCeil) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	num, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathCeil: argument must be a number")
	}
	return math.Ceil(num), nil
}

func (n *NativeMathCeil) String() string { return "<native fn mathCeil>" }

type NativeMathRound struct{}

func (n *NativeMathRound) Arity() int { return 1 }

func (n *NativeMathRound) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	num, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathRound: argument must be a number")
	}
	return math.Round(num), nil
}

func (n *NativeMathRound) String() string { return "<native fn mathRound>" }

type NativeMathAbs struct{}

func (n *NativeMathAbs) Arity() int { return 1 }

func (n *NativeMathAbs) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	num, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathAbs: argument must be a number")
	}
	return math.Abs(num), nil
}

func (n *NativeMathAbs) String() string { return "<native fn mathAbs>" }

type NativeMathMax struct{}

func (n *NativeMathMax) Arity() int { return 2 }

func (n *NativeMathMax) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	a, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathMax: first argument must be a number")
	}
	b, ok := arguments[1].(float64)
	if !ok {
		return nil, fmt.Errorf("mathMax: second argument must be a number")
	}
	return math.Max(a, b), nil
}

func (n *NativeMathMax) String() string { return "<native fn mathMax>" }

type NativeMathMin struct{}

func (n *NativeMathMin) Arity() int { return 2 }

func (n *NativeMathMin) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	a, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathMin: first argument must be a number")
	}
	b, ok := arguments[1].(float64)
	if !ok {
		return nil, fmt.Errorf("mathMin: second argument must be a number")
	}
	return math.Min(a, b), nil
}

func (n *NativeMathMin) String() string { return "<native fn mathMin>" }

type NativeMathRandom struct{}

func (n *NativeMathRandom) Arity() int { return 0 }

func (n *NativeMathRandom) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return rand.Float64(), nil
}

func (n *NativeMathRandom) String() string { return "<native fn mathRandom>" }

type NativeMathPow struct{}

func (n *NativeMathPow) Arity() int { return 2 }

func (n *NativeMathPow) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	base, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathPow: first argument must be a number")
	}
	exp, ok := arguments[1].(float64)
	if !ok {
		return nil, fmt.Errorf("mathPow: second argument must be a number")
	}
	return math.Pow(base, exp), nil
}

func (n *NativeMathPow) String() string { return "<native fn mathPow>" }

type NativeMathSqrt struct{}

func (n *NativeMathSqrt) Arity() int { return 1 }

func (n *NativeMathSqrt) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	num, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("mathSqrt: argument must be a number")
	}
	return math.Sqrt(num), nil
}

func (n *NativeMathSqrt) String() string { return "<native fn mathSqrt>" }

// Type Conversion Functions

type NativeToNumber struct{}

func (n *NativeToNumber) Arity() int { return 1 }

func (n *NativeToNumber) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	switch v := arguments[0].(type) {
	case float64:
		return v, nil
	case string:
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("toNumber: cannot convert string to number")
		}
		return num, nil
	case bool:
		if v {
			return float64(1), nil
		}
		return float64(0), nil
	default:
		return nil, fmt.Errorf("toNumber: cannot convert value to number")
	}
}

func (n *NativeToNumber) String() string { return "<native fn toNumber>" }

type NativeToString struct{}

func (n *NativeToString) Arity() int { return 1 }

func (n *NativeToString) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return stringify(arguments[0]), nil
}

func (n *NativeToString) String() string { return "<native fn toString>" }

type NativeToBoolean struct{}

func (n *NativeToBoolean) Arity() int { return 1 }

func (n *NativeToBoolean) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return isTruthy(arguments[0]), nil
}

func (n *NativeToBoolean) String() string { return "<native fn toBoolean>" }

// Utility function for getting current time
type NativeTime struct{}

func (n *NativeTime) Arity() int { return 0 }

func (n *NativeTime) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func (n *NativeTime) String() string { return "<native fn time>" }

// Sleep function
type NativeSleep struct{}

func (n *NativeSleep) Arity() int { return 1 }

func (n *NativeSleep) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	ms, ok := arguments[0].(float64)
	if !ok {
		return nil, fmt.Errorf("sleep: argument must be a number (milliseconds)")
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return nil, nil
}

func (n *NativeSleep) String() string { return "<native fn sleep>" }

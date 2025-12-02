package main

import (
	"fmt"
	"strconv"
	"strings"
)

// --- Callable Interface and Function Type ---

// Interface for anything callable (functions, maybe classes later)
type LangoCallable interface {
	Arity() int                      // Number of arguments expected
	Call(*Interpreter, []interface{}) (interface{}, error) // Execution logic
	String() string                  // How it represents itself as a string
}

// Represents a user-defined function OR method
type LangoFunction struct {
	declaration *Function    // The Function Stmt node
	closure     *Environment // Environment where the function/method was declared or bound
	// TODO: Add isInitializer flag if implementing init methods
}

// Arity returns the number of parameters
func (f *LangoFunction) Arity() int {
	return len(f.declaration.Params)
}

// Bind creates a new function closure with 'this' defined.
func (f *LangoFunction) Bind(instance *LangoInstance) *LangoFunction {
	// Create an environment nested inside the method's original closure.
	environment := NewEnvironment(f.closure)
	// Define 'this' in that new environment.
	environment.Define("this", instance)
	// Return a new LangoFunction with the updated closure.
	// The declaration itself is reused.
	return &LangoFunction{declaration: f.declaration, closure: environment}
}

// Call executes the function
func (f *LangoFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	// Environment for execution is now the function's closure (which might already contain 'this')
	environment := NewEnvironment(f.closure)

	// Bind arguments to parameter names
	for i, param := range f.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	// Execute with panic/recover for returns
	var returnValue interface{} = nil
	err := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				if retVal, ok := r.(ReturnValue); ok {
					returnValue = retVal.value
					err = nil
				} else {
					panic(r)
				}
			}
		}()
		// Execute body in the environment created *for this call* (which encloses the closure)
		err = interpreter.executeBlock(f.declaration.Body, environment)
		return err
	}()

	if err != nil {
		return nil, err
	}

	return returnValue, nil
}

// String representation of the function
func (f *LangoFunction) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

// --- Return Value Handling ---

// Custom type to signal a return value via panic/recover
type ReturnValue struct {
	value interface{}
}

// --- Class and Instance Types ---

// Represents a Lango class
type LangoClass struct {
	name    string                // Class name
	methods map[string]*LangoFunction // Methods defined in the class
}

// Arity for a class call (constructor) - typically 0 unless we add initializers
func (c *LangoClass) Arity() int {
	// TODO: Implement initializer arity if needed
	return 0
}

// Call for a class creates an instance (acts as constructor)
func (c *LangoClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := NewLangoInstance(c)
	// TODO: Add initializer logic here if implementing init methods
	return instance, nil
}

// String representation of the class
func (c *LangoClass) String() string {
	return fmt.Sprintf("<class %s>", c.name)
}

// Represents an instance of a Lango class
type LangoInstance struct {
	class  *LangoClass          // The instance's class
	fields map[string]interface{} // Instance fields (properties)
}

func NewLangoInstance(class *LangoClass) *LangoInstance {
	return &LangoInstance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

// Get a property or method from the instance
func (i *LangoInstance) Get(name Token) (interface{}, error) {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value, nil
	}

	if method, ok := i.class.methods[name.Lexeme]; ok {
		// Bind 'this' to the method before returning it
		return method.Bind(i), nil
	}

	return nil, fmt.Errorf("[line %d] Undefined property '%s'.", name.Line, name.Lexeme)
}

// Set a property on the instance
func (i *LangoInstance) Set(name Token, value interface{}) {
	i.fields[name.Lexeme] = value
}

// String representation of an instance
func (i *LangoInstance) String() string {
	return fmt.Sprintf("<instance of %s>", i.class.name)
}

// --- Interpreter Struct and Methods ---

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	env := NewEnvironment(nil)
	
	// Register HTTP functions
	env.Define("httpGet", &NativeHttpGet{})
	env.Define("httpPost", &NativeHttpPost{})
	env.Define("httpPut", &NativeHttpPut{})
	env.Define("httpDelete", &NativeHttpDelete{})
	
	// Register JSON functions
	env.Define("jsonParse", &NativeJsonParse{})
	env.Define("jsonStringify", &NativeJsonStringify{})
	env.Define("jsonStringifyPretty", &NativeJsonStringifyPretty{})
	
	// Register file I/O functions
	env.Define("fileRead", &NativeFileRead{})
	env.Define("fileWrite", &NativeFileWrite{})
	env.Define("fileAppend", &NativeFileAppend{})
	env.Define("fileExists", &NativeFileExists{})
	env.Define("fileDelete", &NativeFileDelete{})
	env.Define("fileList", &NativeFileList{})
	env.Define("fileCopy", &NativeFileCopy{})
	
	// Register string functions
	env.Define("strLen", &NativeStrLen{})
	env.Define("strUpper", &NativeStrUpper{})
	env.Define("strLower", &NativeStrLower{})
	env.Define("strSplit", &NativeStrSplit{})
	env.Define("strJoin", &NativeStrJoin{})
	env.Define("strContains", &NativeStrContains{})
	env.Define("strReplace", &NativeStrReplace{})
	env.Define("strTrim", &NativeStrTrim{})
	
	// Register array functions
	env.Define("arrayLen", &NativeArrayLen{})
	env.Define("arrayPush", &NativeArrayPush{})
	env.Define("arrayPop", &NativeArrayPop{})
	env.Define("arraySlice", &NativeArraySlice{})
	env.Define("arrayConcat", &NativeArrayConcat{})
	
	// Register math functions
	env.Define("mathFloor", &NativeMathFloor{})
	env.Define("mathCeil", &NativeMathCeil{})
	env.Define("mathRound", &NativeMathRound{})
	env.Define("mathAbs", &NativeMathAbs{})
	env.Define("mathMax", &NativeMathMax{})
	env.Define("mathMin", &NativeMathMin{})
	env.Define("mathRandom", &NativeMathRandom{})
	env.Define("mathPow", &NativeMathPow{})
	env.Define("mathSqrt", &NativeMathSqrt{})
	
	// Register type conversion functions
	env.Define("toNumber", &NativeToNumber{})
	env.Define("toString", &NativeToString{})
	env.Define("toBoolean", &NativeToBoolean{})
	
	// Register utility functions
	env.Define("time", &NativeTime{})
	env.Define("sleep", &NativeSleep{})
	
	return &Interpreter{
		environment: env,
	}
}


func (i *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			runtimeError(err)
		}
	}
}

func (i *Interpreter) execute(stmt Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBlockStmt(stmt *Block) (interface{}, error) {
	i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
	return nil, nil
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := i.environment
	defer func() { i.environment = previous }()
	i.environment = environment
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *Expression) (interface{}, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitForStmt(stmt *For) (interface{}, error) {
	if stmt.Initializer != nil {
		_, err := i.execute(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	for {
		if stmt.Condition != nil {
			cond, err := i.evaluate(stmt.Condition)
			if err != nil {
				return nil, err
			}
			if !i.isTruthy(cond) {
				break
			}
		}

		_, err := i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}

		if stmt.Increment != nil {
			_, err := i.evaluate(stmt.Increment)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt *If) (interface{}, error) {
	cond, err := i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(cond) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) (interface{}, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(i.stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case MINUS:
		if value, ok := right.(float64); ok {
			return -value, nil
		}
		return nil, i.error(*expr.Operator, "Operand must be a number.")
	case BANG:
		return !i.isTruthy(right), nil
	}

	return nil, i.error(*expr.Operator, "Unexpected unary operator.")
}

func (i *Interpreter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return i.environment.Get(*expr.Name)
}

func (i *Interpreter) VisitVarStmt(stmt *Var) (interface{}, error) {
	if _, ok := i.environment.values[stmt.Name.Lexeme]; !ok {
		i.environment.Define(stmt.Name.Lexeme, nil)
	}

	if stmt.Initializer != nil {
		value, err := i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
		i.environment.Assign(*stmt.Name, value)
	}
	return nil, nil
}

func (i *Interpreter) VisitWhileStmt(stmt *While) (interface{}, error) {
	for {
		cond, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
		if !i.isTruthy(cond) {
			break
		}
		_, err = i.execute(stmt.Body)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case MINUS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l - r, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be numbers.")
	case SLASH:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					return nil, i.error(*expr.Operator, "Division by zero.")
				}
				return l / r, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be numbers.")
	case STAR:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l * r, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be numbers.")
	case MOD:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				if r == 0 {
					return nil, i.error(*expr.Operator, "Modulo by zero.")
				}
				return float64(int(l) % int(r)), nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands of modulo must be numbers.")
	case PLUS:
		if lNum, ok := left.(float64); ok {
			if rNum, ok := right.(float64); ok {
				return lNum + rNum, nil
			}
		}
		if lStr, ok := left.(string); ok {
			if rStr, ok := right.(string); ok {
				return lStr + rStr, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be two numbers or two strings.")
	case GREATER:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l > r, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be numbers.")
	case GREATER_EQUAL:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l >= r, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be numbers.")
	case LESS:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l < r, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be numbers.")
	case LESS_EQUAL:
		if l, ok := left.(float64); ok {
			if r, ok := right.(float64); ok {
				return l <= r, nil
			}
		}
		return nil, i.error(*expr.Operator, "Operands must be numbers.")
	case BANG_EQUAL:
		return !i.isEqual(left, right), nil
	case EQUAL_EQUAL:
		return i.isEqual(left, right), nil
	}

	return nil, i.error(*expr.Operator, "Unexpected binary operator.")
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	err = i.environment.Assign(*expr.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr *Logical) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == OR {
		if i.isTruthy(left) {
			return true, nil
		}
	} else {
		if !i.isTruthy(left) {
			return false, nil
		}
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	return i.isTruthy(right), nil
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}
	if f, ok := object.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	// Add handling for LangoCallable (specifically LangoFunction)
	if callable, ok := object.(LangoCallable); ok {
		return callable.String()
	}
	// Add handling for LangoInstance
	if instance, ok := object.(*LangoInstance); ok {
		return instance.String()
	}
	// Add handling for slices (arrays)
	if array, ok := object.([]interface{}); ok {
		var builder strings.Builder
		builder.WriteString("[")
		for idx, val := range array {
			if idx > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(i.stringify(val))
		}
		builder.WriteString("]")
		return builder.String()
	}
	return fmt.Sprintf("%v", object)
}

func (i *Interpreter) error(token Token, message string) error {
	err := fmt.Errorf("[line %d] Runtime error at '%s': %s", token.Line, token.Lexeme, message)
	runtimeError(err)
	return err
}

// Add implementation for Function statement
func (i *Interpreter) VisitFunctionStmt(stmt *Function) (interface{}, error) {
	// Create the function object, capturing the current environment as the closure
	function := &LangoFunction{declaration: stmt, closure: i.environment}
	// Define the function in the current environment
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil // Function definition doesn't produce a value
}

// Add implementation for Call expression
func (i *Interpreter) VisitCallExpr(expr *Call) (interface{}, error) {
	// 1. Evaluate the Callee expression
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	// 2. Evaluate the Arguments
	arguments := []interface{}{}
	for _, argExpr := range expr.Arguments {
		argValue, err := i.evaluate(argExpr)
		if err != nil {
			return nil, err // Error during argument evaluation
		}
		arguments = append(arguments, argValue)
	}

	// 3. Check if the Callee is actually callable
	function, ok := callee.(LangoCallable)
	if !ok {
		// Use the call expression's parenthesis token for error location
		return nil, i.error(expr.Paren, "Can only call functions and classes.")
	}

	// 4. Check Arity (number of arguments)
	if len(arguments) != function.Arity() {
		return nil, i.error(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.",
			function.Arity(), len(arguments)))
	}

	// 5. Perform the call
	return function.Call(i, arguments)
}

// Add implementation for Return statement
func (i *Interpreter) VisitReturnStmt(stmt *Return) (interface{}, error) {
	var value interface{} = nil
	var err error
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err // Error evaluating the return value
		}
	}

	// Panic with a special value to unwind the stack
	panic(ReturnValue{value: value})
	// This return is unreachable but required by Go
	// return nil, nil
}

// Add implementation for Class statement
func (i *Interpreter) VisitClassStmt(stmt *Class) (interface{}, error) {
	i.environment.Define(stmt.Name.Lexeme, nil) // Define class name early for self-references in methods

	methods := make(map[string]*LangoFunction)
	for _, methodStmt := range stmt.Methods {
		// Create LangoFunction for each method, capturing the current environment
		// Note: The closure here doesn't automatically include 'this'. Binding happens during method lookup/call.
		function := &LangoFunction{declaration: methodStmt, closure: i.environment}
		methods[methodStmt.Name.Lexeme] = function
	}

	// Create the LangoClass object
	class := &LangoClass{name: stmt.Name.Lexeme, methods: methods}

	// Assign the fully formed class object to the name defined earlier
	err := i.environment.Assign(stmt.Name, class)
	if err != nil { return nil, err } 

	return nil, nil
}

// Add implementation for Get expression
func (i *Interpreter) VisitGetExpr(expr *Get) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	// Check if the object is an instance
	if instance, ok := object.(*LangoInstance); ok {
		// Call the instance's Get method
		// Need to pass the Token for error reporting inside Get
		return instance.Get(expr.Name)
	}
	
	// Check if the object is a map (from JSON parsing)
	if mapObj, ok := object.(map[string]interface{}); ok {
		value, exists := mapObj[expr.Name.Lexeme]
		if !exists {
			return nil, i.error(expr.Name, fmt.Sprintf("Undefined property '%s'.", expr.Name.Lexeme))
		}
		return value, nil
	}
	
	// If not an instance or map, it's an error
	return nil, i.error(expr.Name, "Only instances and objects have properties.")
}

// Add implementation for Set expression
func (i *Interpreter) VisitSetExpr(expr *Set) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	// Check if the object is an instance
	instance, ok := object.(*LangoInstance)
	if !ok {
		return nil, i.error(expr.Name, "Only instances have fields.")
	}

	// Evaluate the value being assigned
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	// Call the instance's Set method
	instance.Set(expr.Name, value)

	// Assignment expression returns the assigned value
	return value, nil 
}

// Add implementation for This expression
func (i *Interpreter) VisitThisExpr(expr *This) (interface{}, error) {
	// Look up the special variable "this" in the environment
	// We rely on the environment chain created during method binding (in LangoFunction.Bind)
	return i.environment.LookUp("this", expr.Keyword) 
}

// Add implementation for ArrayLiteral expression
func (i *Interpreter) VisitArrayLiteralExpr(expr *ArrayLiteral) (interface{}, error) {
	values := make([]interface{}, len(expr.Values))
	for idx, valExpr := range expr.Values {
		val, err := i.evaluate(valExpr)
		if err != nil {
			return nil, err
		}
		values[idx] = val
	}
	// Return the Go slice directly
	return values, nil 
}

// Add implementation for ArrayIndex expression
func (i *Interpreter) VisitArrayIndexExpr(expr *ArrayIndex) (interface{}, error) {
	// Evaluate the object being indexed
	arrayObj, err := i.evaluate(expr.Array)
	if err != nil {
		return nil, err
	}

	// Check if it's actually a slice
	arrayVal, ok := arrayObj.([]interface{})
	if !ok {
		return nil, i.error(expr.Bracket, "Can only index arrays (lists).")
	}

	// Evaluate the index
	indexObj, err := i.evaluate(expr.Index)
	if err != nil {
		return nil, err
	}

	// Check if the index is a number
	indexNum, ok := indexObj.(float64)
	if !ok {
		return nil, i.error(expr.Bracket, "Array index must be a number.")
	}

	// Check if the index is an integer
	indexInt := int(indexNum)
	if float64(indexInt) != indexNum {
		return nil, i.error(expr.Bracket, "Array index must be an integer.")
	}

	// Check bounds
	if indexInt < 0 || indexInt >= len(arrayVal) {
		return nil, i.error(expr.Bracket, fmt.Sprintf("Array index out of bounds (%d for size %d).", indexInt, len(arrayVal)))
	}

	// Return the element
	return arrayVal[indexInt], nil
}

// Add implementation for ArrayAssign expression
func (i *Interpreter) VisitArrayAssignExpr(expr *ArrayAssign) (interface{}, error) {
	// Evaluate the object being indexed
	arrayObj, err := i.evaluate(expr.Assignee.Array)
	if err != nil {
		return nil, err
	}

	// Check if it's a slice
	arrayVal, ok := arrayObj.([]interface{})
	if !ok {
		return nil, i.error(expr.Assignee.Bracket, "Can only index arrays (lists).")
	}

	// Evaluate the index
	indexObj, err := i.evaluate(expr.Assignee.Index)
	if err != nil {
		return nil, err
	}

	// Check if the index is a number
	indexNum, ok := indexObj.(float64)
	if !ok {
		return nil, i.error(expr.Assignee.Bracket, "Array index must be a number.")
	}

	// Check if the index is an integer
	indexInt := int(indexNum)
	if float64(indexInt) != indexNum {
		return nil, i.error(expr.Assignee.Bracket, "Array index must be an integer.")
	}

	// Check bounds
	if indexInt < 0 || indexInt >= len(arrayVal) {
		return nil, i.error(expr.Assignee.Bracket, fmt.Sprintf("Array index out of bounds (%d for size %d).", indexInt, len(arrayVal)))
	}

	// Evaluate the value to assign
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	// Perform the assignment
	arrayVal[indexInt] = value

	// Assignment returns the assigned value
	return value, nil
}

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// NativeFileRead reads a file and returns its content as a string
type NativeFileRead struct{}

func (n *NativeFileRead) Arity() int {
	return 1
}

func (n *NativeFileRead) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("fileRead requires exactly 1 argument")
	}

	path, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("fileRead: path must be a string")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("fileRead: failed to read file - %v", err)
	}

	return string(content), nil
}

func (n *NativeFileRead) String() string {
	return "<native fn fileRead>"
}

// NativeFileWrite writes content to a file (overwrites if exists)
type NativeFileWrite struct{}

func (n *NativeFileWrite) Arity() int {
	return 2
}

func (n *NativeFileWrite) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 2 {
		return nil, fmt.Errorf("fileWrite requires exactly 2 arguments")
	}

	path, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("fileWrite: path must be a string")
	}

	content, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("fileWrite: content must be a string")
	}

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return nil, fmt.Errorf("fileWrite: failed to write file - %v", err)
	}

	return nil, nil
}

func (n *NativeFileWrite) String() string {
	return "<native fn fileWrite>"
}

// NativeFileAppend appends content to a file
type NativeFileAppend struct{}

func (n *NativeFileAppend) Arity() int {
	return 2
}

func (n *NativeFileAppend) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 2 {
		return nil, fmt.Errorf("fileAppend requires exactly 2 arguments")
	}

	path, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("fileAppend: path must be a string")
	}

	content, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("fileAppend: content must be a string")
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("fileAppend: failed to open file - %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return nil, fmt.Errorf("fileAppend: failed to append to file - %v", err)
	}

	return nil, nil
}

func (n *NativeFileAppend) String() string {
	return "<native fn fileAppend>"
}

// NativeFileExists checks if a file exists
type NativeFileExists struct{}

func (n *NativeFileExists) Arity() int {
	return 1
}

func (n *NativeFileExists) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("fileExists requires exactly 1 argument")
	}

	path, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("fileExists: path must be a string")
	}

	_, err := os.Stat(path)
	return err == nil, nil
}

func (n *NativeFileExists) String() string {
	return "<native fn fileExists>"
}

// NativeFileDelete deletes a file
type NativeFileDelete struct{}

func (n *NativeFileDelete) Arity() int {
	return 1
}

func (n *NativeFileDelete) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("fileDelete requires exactly 1 argument")
	}

	path, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("fileDelete: path must be a string")
	}

	err := os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("fileDelete: failed to delete file - %v", err)
	}

	return nil, nil
}

func (n *NativeFileDelete) String() string {
	return "<native fn fileDelete>"
}

// NativeFileList lists files in a directory
type NativeFileList struct{}

func (n *NativeFileList) Arity() int {
	return 1
}

func (n *NativeFileList) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 1 {
		return nil, fmt.Errorf("fileList requires exactly 1 argument")
	}

	path, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("fileList: path must be a string")
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("fileList: failed to list directory - %v", err)
	}

	result := make([]interface{}, 0, len(files))
	for _, file := range files {
		result = append(result, file.Name())
	}

	return result, nil
}

func (n *NativeFileList) String() string {
	return "<native fn fileList>"
}

// NativeFileCopy copies a file
type NativeFileCopy struct{}

func (n *NativeFileCopy) Arity() int {
	return 2
}

func (n *NativeFileCopy) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) != 2 {
		return nil, fmt.Errorf("fileCopy requires exactly 2 arguments")
	}

	srcPath, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("fileCopy: source path must be a string")
	}

	dstPath, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("fileCopy: destination path must be a string")
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("fileCopy: failed to open source file - %v", err)
	}
	defer srcFile.Close()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return nil, fmt.Errorf("fileCopy: failed to create destination directory - %v", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("fileCopy: failed to create destination file - %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return nil, fmt.Errorf("fileCopy: failed to copy file - %v", err)
	}

	return nil, nil
}

func (n *NativeFileCopy) String() string {
	return "<native fn fileCopy>"
}

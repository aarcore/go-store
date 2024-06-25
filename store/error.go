package store

import "errors"

var (
	FileNotFoundError = errors.New("File not found")
	FileNoPermissionError = errors.New("Permission denied")
)
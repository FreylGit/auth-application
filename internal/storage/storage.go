package storage

import "errors"

var (
	ErrNotFound    = errors.New("is not found")
	ErrorNoUnique  = errors.New("value is not unique")
	ErrorSave      = errors.New("failed to save")
	ErrorUpdate    = errors.New("failed to update")
	ErrorSqlSyntax = errors.New("sql syntax is not correct")
	ErrorScan      = errors.New("failed parse query result")
)

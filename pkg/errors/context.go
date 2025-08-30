package errors

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// BasicContext provides minimal context for certificate errors
type BasicContext struct {
	Component string    `json:"component"`
	Operation string    `json:"operation"`
	OS        string    `json:"os"`
	Timestamp time.Time `json:"timestamp"`
	User      string    `json:"user"`
}

// CaptureBasicContext collects minimal error context
func CaptureBasicContext(component, operation string) *BasicContext {
	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}
	
	return &BasicContext{
		Component: component,
		Operation: operation,
		OS:        runtime.GOOS,
		Timestamp: time.Now(),
		User:      user,
	}
}

// ToString provides a string representation of the context
func (ctx *BasicContext) ToString() string {
	return fmt.Sprintf("Component: %s, Operation: %s, OS: %s, User: %s, Time: %s",
		ctx.Component, ctx.Operation, ctx.OS, ctx.User, ctx.Timestamp.Format(time.RFC3339))
}

// EnrichError adds basic context to a CertificateError
func EnrichError(err *CertificateError, component, operation string) *CertificateError {
	ctx := CaptureBasicContext(component, operation)
	
	if err.Details == nil {
		err.Details = make(map[string]string)
	}
	
	err.Details["component"] = ctx.Component
	err.Details["operation"] = ctx.Operation
	err.Details["os"] = ctx.OS
	err.Details["user"] = ctx.User
	err.Details["timestamp"] = ctx.Timestamp.Format(time.RFC3339)
	
	return err
}
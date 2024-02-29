package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

const (
	// Info return "info"
	Info LogLevel = "info"
	// Warning return "warning"
	Warning LogLevel = "warning"
	// Error return "error"
	Error LogLevel = "error"
	// Fatal return "fatal"
	Fatal LogLevel = "fatal"
)

type LogLevel string

type logger struct {
	ServiceName string
}

func NewLogger(serviceName string) logger {
	return logger{
		ServiceName: serviceName,
	}
}

func (logger *logger) Log(ctx context.Context, eventName string, level LogLevel, attributes map[string]interface{}) {
	output := map[string]interface{}{
		"event":      eventName,
		"level":      level,
		"service":    logger.ServiceName,
		"properties": attributes,
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		log.Println("error marshaling data:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (logger *logger) Info(ctx context.Context, eventName string, attributes map[string]interface{}) {
	logger.Log(ctx, eventName, Info, attributes)
}

func (logger *logger) Warning(ctx context.Context, eventName string, attributes map[string]interface{}) {
	logger.Log(ctx, eventName, Warning, attributes)
}

func (logger *logger) Error(ctx context.Context, eventName string, attributes map[string]interface{}, err error) {
	attributes["error"] = err.Error()

	logger.Log(ctx, eventName, Error, attributes)
}

func (logger *logger) Fatal(ctx context.Context, eventName string, attributes map[string]interface{}) {
	logger.Log(ctx, eventName, Fatal, attributes)
}

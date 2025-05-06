package model

import (
	"biathlon-system-prototype/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestParseEvent_ValidInput проверяет корректный разбор строки события без дополнительных параметров
func TestParseEvent_ValidInput(t *testing.T) {
	input := "[12:30:00.000] 1 42"
	event, err := ParseEvent(input)
	require.NoError(t, err)
	assert.Equal(t, Registered, event.Id)
	assert.Equal(t, 42, event.CompetitorID)
	assert.Equal(t, "12:30:00.000", event.Time.Format(config.TimeLayout))
}

// TestParseEvent_WithExtraParams проверяет разбор события с дополнительными параметрами
func TestParseEvent_WithExtraParams(t *testing.T) {
	input := "[12:35:10.000] 6 7 1"
	event, err := ParseEvent(input)
	require.NoError(t, err)
	assert.Equal(t, TargetHit, event.Id)
	assert.Equal(t, 7, event.CompetitorID)
	assert.Equal(t, []string{"1"}, event.ExtraParams)
}

// TestParseEvent_InvalidFormat проверяет обработку строки с некорректным форматом без скобок
func TestParseEvent_InvalidFormat(t *testing.T) {
	input := "12:35:10 1 7"
	_, err := ParseEvent(input)
	assert.Error(t, err)
}

// TestParseEvent_BadTime проверяет обработку строки с некорректным временем
func TestParseEvent_BadTime(t *testing.T) {
	input := "[badtime] 1 7"
	_, err := ParseEvent(input)
	assert.Error(t, err)
}

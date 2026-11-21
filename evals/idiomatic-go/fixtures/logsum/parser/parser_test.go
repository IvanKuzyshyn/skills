package parser

import "testing"

var cfg = ParserConfig{TimeLayout: "2006-01-02T15:04:05Z07:00"}

func TestParseLineInfo(t *testing.T) {
	e, err := ParseLine(cfg, "2024-01-01T00:00:00Z INFO r1 started")
	if err != nil {
		t.Fatal(err)
	}
	if e.Level != "INFO" {
		t.Fatal("wrong level")
	}
}

func TestParseLineError(t *testing.T) {
	e, err := ParseLine(cfg, "2024-01-01T00:00:00Z ERROR r2 boom")
	if err != nil {
		t.Fatal(err)
	}
	if e.Level != "ERROR" {
		t.Fatal("wrong level")
	}
}

func TestParseLineBad(t *testing.T) {
	_, err := ParseLine(cfg, "garbage")
	if err == nil {
		t.Fatal("expected error")
	}
}

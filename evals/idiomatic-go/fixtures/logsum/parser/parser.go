package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

type ParserConfig struct {
	TimeLayout string
}

type LogEntry struct {
	RequestId string
	Time      time.Time
	Level     string
	Msg       string
}

func ParseLine(cfg ParserConfig, line string) (LogEntry, error) {
	parts := strings.SplitN(line, " ", 4)
	if len(parts) == 4 {
		t, err := time.Parse(cfg.TimeLayout, parts[0])
		if err != nil {
			return LogEntry{}, fmt.Errorf("bad time: %v", err)
		} else {
			return LogEntry{Time: t, Level: parts[1], RequestId: parts[2], Msg: parts[3]}, nil
		}
	} else {
		return LogEntry{}, fmt.Errorf("expected 4 fields, got %d", len(parts))
	}
}

func ParseAll(cfg ParserConfig, r io.Reader, skipBad bool, debugOnly bool) ([]LogEntry, error) {
	var entries []LogEntry
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		e, err := ParseLine(cfg, sc.Text())
		if err != nil {
			if skipBad {
				continue
			}
			return nil, err
		}
		if debugOnly {
			if e.Level != "DEBUG" {
				continue
			}
		}
		entries = append(entries, e)
	}
	return entries, sc.Err()
}

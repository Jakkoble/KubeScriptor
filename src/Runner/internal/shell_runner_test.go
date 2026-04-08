package internal_test

import (
	"context"
	"testing"

	"github.com/Jakkoble/HexaTask/src/Runner/internal"
)

type captureStreamer struct {
	logs    []string
	stderrs []bool
}

func (c *captureStreamer) SendLog(msg string, isStderr bool) error {
	c.logs = append(c.logs, msg)
	c.stderrs = append(c.stderrs, isStderr)
	return nil
}
func (c *captureStreamer) Close() error { return nil }

func TestShellRunner_Echo(t *testing.T) {
	streamer := &captureStreamer{}
	code, err := (&internal.ShellCommandRunner{}).Run(context.Background(), "echo hello", streamer)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}

	if len(streamer.logs) == 0 || streamer.logs[0] != "hello" {
		t.Errorf("expected log 'hello', got %v", streamer.logs)
	}
}

func TestShellRunner_Stderr(t *testing.T) {
	streamer := &captureStreamer{}
	(&internal.ShellCommandRunner{}).Run(context.Background(), "echo errline >&2", streamer)

	for i, isStderr := range streamer.stderrs {
		if isStderr {
			if streamer.logs[i] != "errline" {
				t.Errorf("expected stderr log 'errline', got '%s'", streamer.logs[i])
			}
			return
		}
	}
	t.Error("expected at least one log with isStderr=true")
}

func TestShellRunner_ExitCode(t *testing.T) {
	code, err := (&internal.ShellCommandRunner{}).Run(context.Background(), "exit 42", &captureStreamer{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 42 {
		t.Errorf("expected exit code 42, got %d", code)
	}
}

func TestShellRunner_InvalidCommand(t *testing.T) {
	_, err := (&internal.ShellCommandRunner{}).Run(context.Background(), "this_command_does_not_exist_xyz", &captureStreamer{})

	if err != nil {
		t.Errorf("expected nil error (shell handles unknown commands), got %v", err)
	}
}

func TestShellRunner_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	code, _ := (&internal.ShellCommandRunner{}).Run(ctx, "sleep 10", &captureStreamer{})

	if code == 0 {
		t.Error("expected non-zero exit code for cancelled context")
	}
}

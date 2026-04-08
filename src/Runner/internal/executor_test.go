package internal_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Jakkoble/HexaTask/src/Runner/internal"
)

type mockStreamer struct {
	closed bool
}

func (m *mockStreamer) SendLog(_ string, _ bool) error { return nil }
func (m *mockStreamer) Close() error {
	m.closed = true
	return nil
}

type mockRunner struct {
	exitCodes map[string]int
	runError  map[string]error
	executed  []string
}

func (m *mockRunner) Run(_ context.Context, cmdStr string, _ internal.ILogStreamer) (int, error) {
	m.executed = append(m.executed, cmdStr)

	if err, ok := m.runError[cmdStr]; ok {
		return 1, err
	}

	if code, ok := m.exitCodes[cmdStr]; ok {
		return code, nil
	}

	return 0, nil
}

func TestExecute_AllSucceed(t *testing.T) {
	runner := &mockRunner{}
	streamer := &mockStreamer{}

	code := internal.NewTaskExecutor(runner, streamer).
		Execute(context.Background(), []string{"a", "b", "c"})

	if code != 0 {
		t.Errorf("expected 0, got %d", code)
	}
}

func TestExecute_FirstCommandFails(t *testing.T) {
	runner := &mockRunner{
		exitCodes: map[string]int{"fail": 1},
	}
	streamer := &mockStreamer{}

	internal.NewTaskExecutor(runner, streamer).
		Execute(context.Background(), []string{"fail", "second"})

	if len(runner.executed) != 1 {
		t.Errorf("expected 1 executed command, got %d", len(runner.executed))
	}
}

func TestExecute_CommandFails(t *testing.T) {
	runner := &mockRunner{
		exitCodes: map[string]int{"cmd": 42},
	}
	streamer := &mockStreamer{}

	code := internal.NewTaskExecutor(runner, streamer).
		Execute(context.Background(), []string{"cmd"})

	if code != 42 {
		t.Errorf("expected exit code 42, got %d", code)
	}
}

func TestExecute_RunnerReturnsError(t *testing.T) {
	runner := &mockRunner{
		runError: map[string]error{"cmd": errors.New("unexpected failure")},
	}
	streamer := &mockStreamer{}

	code := internal.NewTaskExecutor(runner, streamer).
		Execute(context.Background(), []string{"cmd"})

	if code != 1 {
		t.Errorf("expected 1 on runner error, got %d", code)
	}
}

func TestExecute_StreamerAlwaysClosed(t *testing.T) {
	runner := &mockRunner{exitCodes: map[string]int{"fail": 1}}
	streamer := &mockStreamer{}

	internal.NewTaskExecutor(runner, streamer).
		Execute(context.Background(), []string{"fail"})

	if !streamer.closed {
		t.Error("expected streamer.Close() to be called even after failure")
	}
}

func TestExecute_EmptyCommands(t *testing.T) {
	runner := &mockRunner{}
	streamer := &mockStreamer{}

	code := internal.NewTaskExecutor(runner, streamer).
		Execute(context.Background(), []string{})

	if code != 0 {
		t.Errorf("expected 0 for empty command list, got %d", code)
	}
}

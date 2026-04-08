package internal

import (
	"context"
)

type ILogStreamer interface {
	SendLog(msg string, isStderr bool) error
	Close(exitCode int) error
}

type ICommandRunner interface {
	Run(ctx context.Context, cmdStr string, streamer ILogStreamer) (exitCode int, err error)
}

type TaskExecutor struct {
	runner   ICommandRunner
	streamer ILogStreamer
}

func NewTaskExecutor(runner ICommandRunner, streamer ILogStreamer) *TaskExecutor {
	return &TaskExecutor{runner: runner, streamer: streamer}
}

func (e *TaskExecutor) Execute(ctx context.Context, commands []string) int {
	exitCode := 0

	for _, cmdStr := range commands {
		code, err := e.runner.Run(ctx, cmdStr, e.streamer)
		if err != nil || code != 0 {
			exitCode = max(code, 1)
			break
		}
	}

	e.streamer.Close(exitCode)
	return exitCode
}

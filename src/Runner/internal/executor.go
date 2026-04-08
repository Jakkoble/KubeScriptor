package internal

import (
	"context"
)

type ILogStreamer interface {
	SendLog(msg string, isStderr bool) error
	Close() error
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
	defer e.streamer.Close()

	for _, cmdStr := range commands {
		exitCode, err := e.runner.Run(ctx, cmdStr, e.streamer)
		if err != nil {
			return 1
		}

		if exitCode != 0 {
			return exitCode
		}
	}
	return 0
}

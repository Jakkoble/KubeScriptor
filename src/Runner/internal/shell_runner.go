package internal

import (
	"bufio"
	"context"
	"os/exec"
	"sync"
)

type ShellCommandRunner struct{}

func (r *ShellCommandRunner) Run(ctx context.Context, cmdStr string, streamer ILogStreamer) (int, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 1, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 1, err
	}

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			_ = streamer.SendLog(scanner.Text(), false)
		}
	}()

	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			_ = streamer.SendLog(scanner.Text(), true)
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}

		return 1, err
	}

	return 0, nil
}

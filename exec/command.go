package exec

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func RunCommand(ctx context.Context, logger logrus.FieldLogger, command string, cwd string, env []string) error {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	cmd.Env = os.Environ()

	env = append(os.Environ(), env...)
	cmd.Env = env
	cmd.Dir = cwd

	err := cmd.Run()
	if err != nil {
		commandStr := limitString(command, 100)
		if os.Getenv("ATLAS_DEBUG_CMD_FULL_COMMAND") == "true" {
			commandStr = command
		}

		fields := logrus.Fields{
			"command": commandStr,
			"cwd":     cwd,
			"stdout":  outBuf.String(),
			"stderr":  errBuf.String(),
		}

		if os.Getenv("ATLAS_DEBUG_CMD_ENV") == "true" {
			fields["env"] = env
		}

		logger.WithFields(fields).Errorln("Could not run command")
		return fmt.Errorf("could not run command %s: %w", commandStr, err)
	}

	return nil
}

func limitString(str string, to int) string {
	runes := []rune(str)
	if len(runes) > to {
		return string(runes[:to])
	}

	return str
}

func StartCommand(ctx context.Context, logger logrus.FieldLogger, command string, cwd string, env []string) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	cmd.Env = os.Environ()

	env = append(os.Environ(), env...)
	cmd.Env = env
	cmd.Dir = cwd

	err := cmd.Start()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"command": command,
			"cwd":     cwd,
			"env":     env,
			"stdout":  outBuf.String(),
			"stderr":  errBuf.String(),
		}).Errorln("Could not start command")
		return nil, fmt.Errorf("could not run command %s: %w", command, err)
	}

	return cmd, nil
}

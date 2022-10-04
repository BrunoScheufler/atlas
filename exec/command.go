package exec

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

type PrefixWriter struct {
	w      io.Writer
	prefix string
}

func (pw PrefixWriter) Write(p []byte) (n int, err error) {
	lines := bytes.Split(p, []byte{'\n'})
	for _, line := range lines {
		if len(line) > 0 {
			_, err = fmt.Fprintf(pw.w, "%s%s", pw.prefix, line)
			if err != nil {
				return 0, err
			}
		}
	}

	return len(p), nil
}

type RunCommandOptions struct {
	Cwd string
	Env []string

	LogVisible bool
	LogPrefix  string
}

func RunCommand(ctx context.Context, logger logrus.FieldLogger, command string, options RunCommandOptions) error {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	outBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	var wo io.Writer = outBuf
	if options.LogVisible {
		wo = io.MultiWriter(outBuf, PrefixWriter{os.Stdout, options.LogPrefix})
	}

	var we io.Writer = errBuf
	if options.LogVisible {
		we = io.MultiWriter(errBuf, PrefixWriter{os.Stderr, options.LogPrefix})
	}

	cmd.Stdout = wo
	cmd.Stderr = we

	cmd.Env = os.Environ()

	env := append(os.Environ(), options.Env...)
	cmd.Env = env
	cmd.Dir = options.Cwd

	err := cmd.Run()
	if err != nil {
		commandStr := limitString(command, 100)
		if os.Getenv("ATLAS_DEBUG_CMD_FULL_COMMAND") == "true" {
			commandStr = command
		}

		fields := logrus.Fields{
			"command": commandStr,
			"cwd":     options.Cwd,
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

	// Write output to stdout and outBuf
	wo := io.MultiWriter(outBuf, os.Stdout)
	we := io.MultiWriter(errBuf, os.Stderr)

	cmd.Stdout = wo
	cmd.Stderr = we

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

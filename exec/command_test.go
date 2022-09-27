package exec

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRunCommand(t *testing.T) {
	err := RunCommand(context.Background(), logrus.New(), "[ -n \"\"] && echo hello world", "/tmp", nil, false)
	if err != nil {
		t.Error(err)
	}
}

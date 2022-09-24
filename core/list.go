package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/sirupsen/logrus"
)

func List(ctx context.Context, logger logrus.FieldLogger, version, cwd string) error {
	cwd, err := atlasfile.FindRootDir(cwd)
	if err != nil {
		return fmt.Errorf("could not find root directory: %w", err)
	}

	mergedFile, err := atlasfile.Collect(ctx, logger, version, cwd)
	if err != nil {
		return fmt.Errorf("could not collect atlas files: %w", err)
	}

	fmt.Println(mergedFile.String())

	return nil
}

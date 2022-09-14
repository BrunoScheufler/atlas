package atlas

import (
	"context"
	"fmt"
	"github.com/brunoscheufler/atlas/atlasfile"
	"github.com/sirupsen/logrus"
)

func List(ctx context.Context, logger logrus.FieldLogger, cwd string) error {
	mergedFile, err := atlasfile.Collect(ctx, logger, cwd)
	if err != nil {
		return fmt.Errorf("could not collect atlas files: %w", err)
	}

	fmt.Println(mergedFile.String())

	return nil
}

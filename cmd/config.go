package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/marcus-crane/wwenet/config"
	"github.com/marcus-crane/wwenet/storage"
)

func OutputConfig(ctx context.Context, cmd *cli.Command, cfg config.Config, db *storage.Queries) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return nil
	}
	fmt.Printf("%s\n", b)
	return nil
}

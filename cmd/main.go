package main

import (
	"encoding/json"
	"fmt"
	"os"
	imager "imager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	inputFlagName = "input"
)

func main() {
	if _, err := os.Stat(imager.RepoDirectory); err == nil {
		if err := os.RemoveAll(imager.RepoDirectory); err != nil {
			panic(err)
		}
	}

	cmd := &cobra.Command{
		Use:   "imager",
		Short: "Pull docker image data out of your repositories",
		Long: `Imager, a tool to extract docker image information from your repository, useful
				for auditing and validation purposes e.g all images must have tags `,
		RunE: func(cmd *cobra.Command, args []string) error {
			data := map[string]imager.Response{"data": imager.Master(cmd.Flag(inputFlagName).Value.String())}
			content, err := json.MarshalIndent(data, "", "\t")
			if err != nil {
				return err
			}

			fmt.Println(string(content))
			return nil
		},
	}

	cmd.Flags().String(inputFlagName, "", "url to list of repositories to fetch information from")
	if err := cmd.MarkFlagRequired(inputFlagName); err != nil {
		log.WithField("instance", "main").Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.WithField("instance", "main").Fatal(err)
	}
}

package images

import (
	"fmt"

	"github.com/containers/libpod/cmd/podman/registry"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	rmDescription = "Removes one or more previously pulled or locally created images."
	rmCmd         = &cobra.Command{
		Use:   "rm [flags] IMAGE [IMAGE...]",
		Short: "Removes one or more images from local storage",
		Long:  rmDescription,
		RunE:  rm,
		Example: `podman image rm imageID
  podman image rm --force alpine
  podman image rm c4dfb1609ee2 93fd78260bd1 c0ed59d05ff7`,
	}

	imageOpts = entities.ImageRemoveOptions{}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: rmCmd,
		Parent:  imageCmd,
	})

	imageRemoveFlagSet(rmCmd.Flags())
}

func imageRemoveFlagSet(flags *pflag.FlagSet) {
	flags.BoolVarP(&imageOpts.All, "all", "a", false, "Remove all images")
	flags.BoolVarP(&imageOpts.Force, "force", "f", false, "Force Removal of the image")
}

func rm(cmd *cobra.Command, args []string) error {
	if len(args) < 1 && !imageOpts.All {
		return errors.Errorf("image name or ID must be specified")
	}
	if len(args) > 0 && imageOpts.All {
		return errors.Errorf("when using the --all switch, you may not pass any images names or IDs")
	}

	report, err := registry.ImageEngine().Remove(registry.GetContext(), args, imageOpts)
	if report != nil {
		for _, u := range report.Untagged {
			fmt.Println("Untagged: " + u)
		}
		for _, d := range report.Deleted {
			fmt.Println("Deleted: " + d)
		}
		registry.SetExitCode(report.ExitCode)
	}

	return err
}
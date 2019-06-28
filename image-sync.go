package imagesync

import (
	"context"
	"fmt"
	"github.com/containers/image/copy"
	"github.com/containers/image/docker"
	"github.com/containers/image/signature"
	"github.com/containers/image/types"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"
)

func Execute() error {

	app := cli.NewApp()
	app.Name = "imagesync"
	app.Usage = "Sync docker images between repositories."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "src, s",
			Usage: "Reference for the source docker registry.",
		},
		cli.StringFlag{
			Name:  "src-type",
			Usage: "Type of the source docker registry",
			Value: "insecure",
		},
		cli.StringFlag{
			Name:  "dest, d",
			Usage: "Reference for the destination docker registry.",
		},
		cli.StringFlag{
			Name:  "dest-type",
			Usage: "Type of the destination docker registry",
			Value: "insecure",
		},
		cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Use this to copy/override all the tags.",
		},
	}

	app.Action = func(c *cli.Context) error {
		srcRegistry, err := docker.ParseReference(fmt.Sprintf("//%s", c.String("src")))
		if err != nil {
			return errors.WithMessage(err, "parsing source registry url")
		}

		destRegistry, err := docker.ParseReference(fmt.Sprintf("//%s", c.String("dest")))
		if err != nil {
			return errors.WithMessage(err, "parsing destination registry url")
		}

		ctx := context.Background()
		var srcSysCtx *types.SystemContext
		if c.String("src-type") == "insecure" {
			srcSysCtx, err = newSystemContextWithInsecureRegistry()
			if err != nil {
				return errors.WithMessage(err, "setting up source system context")
			}
		}

		var destSysCtx *types.SystemContext
		if c.String("dest-type") == "insecure" {
			destSysCtx, err = newSystemContextWithInsecureRegistry()
			if err != nil {
				return errors.WithMessage(err, "setting up dest system context")
			}
		}

		srcTags, err := docker.GetRepositoryTags(ctx, srcSysCtx, srcRegistry)
		if err != nil {
			return errors.WithMessage(err, "getting source tags")
		}

		destTags, _ := docker.GetRepositoryTags(ctx, destSysCtx, destRegistry)

		targetTags := targetTags(c.Bool("overwrite"), srcTags, destTags)
		if len(targetTags) == 0 {
			fmt.Println("Info: Images in registries are already synced.")
			os.Exit(0)
		}
		fmt.Printf("Info: Starting image sync with %d tags ...\n", len(targetTags))

		for _, tag := range targetTags {
			src, err := docker.ParseReference(fmt.Sprintf("//%s:%s", c.String("src"), tag))
			if err != nil {
				return err
			}

			dest, err := docker.ParseReference(fmt.Sprintf("//%s:%s", c.String("dest"), tag))
			if err != nil {
				return err
			}

			policy, _ := signature.DefaultPolicy(nil)
			policyContext, _ := signature.NewPolicyContext(policy)

			fmt.Printf("Info: Copying image with tag: %s\n", tag)
			_, err = copy.Image(ctx, policyContext, dest, src, &copy.Options{
				ReportWriter:   os.Stdout,
				SourceCtx:      srcSysCtx,
				DestinationCtx: destSysCtx,
			})
			if err != nil {
				return err
			}
		}

		fmt.Println("Info: Registries image sync completed.")

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}
	return nil
}

func targetTags(overwrite bool, src, dest []string) []string {
	if !overwrite {
		return src
	}
	return diff(src, dest)
}

func diff(a, b []string) []string {
	var result []string

	m := make(map[string]bool)
	for _, i := range b {
		m[i] = true
	}

	for _, j := range a {
		if !m[j] {
			result = append(result, j)
		}
	}
	return result
}

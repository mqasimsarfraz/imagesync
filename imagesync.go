package imagesync

import (
	"context"
	"fmt"
	"github.com/containers/image/copy"
	"github.com/containers/image/docker"
	"github.com/containers/image/signature"
	"github.com/containers/image/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"sync"
)

var wg sync.WaitGroup

type copyImageInput struct {
	context           context.Context
	dest              string
	destSystemContext *types.SystemContext
	src               string
	srcSystemContext  *types.SystemContext
}

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
		cli.IntFlag{
			Name:  "max-concurrent-tags",
			Usage: "Maximum number of tags to be synced/copied in parallel.",
			Value: 12,
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
			logrus.Info("Image in registries are already synced")
			os.Exit(0)
		}

		logrus.Infof("Starting image sync with %d tags", len(targetTags))

		// limit the go routines to avoid 429 on registries
		maxConcurrentTags := c.Int("max-concurrent-tags")
		numberOfConcurrentTags := maxConcurrentTags
		if len(targetTags) < maxConcurrentTags {
			numberOfConcurrentTags = len(targetTags)
		}

		ch := make(chan string, len(targetTags))

		wg.Add(numberOfConcurrentTags)

		input := &copyImageInput{
			context:           ctx,
			dest:              c.String("dest"),
			destSystemContext: destSysCtx,
			src:               c.String("src"),
			srcSystemContext:  srcSysCtx,
		}

		for i := 0; i < numberOfConcurrentTags; i++ {
			go func() {
				for {
					tag, ok := <-ch
					if !ok {
						wg.Done()
						return
					}
					input.copyImage(tag)
				}
			}()
		}

		for _, tag := range targetTags {
			ch <- tag
		}

		close(ch)
		wg.Wait()

		logrus.Info("Registries image sync completed.")

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}
	return nil
}

func (ci *copyImageInput) copyImage(tag string) {

	destRef, _ := docker.ParseReference(fmt.Sprintf("//%s:%s", ci.dest, tag))
	srcRef, _ := docker.ParseReference(fmt.Sprintf("//%s:%s", ci.src, tag))

	policy, _ := signature.DefaultPolicy(nil)
	policyContext, _ := signature.NewPolicyContext(policy)

	copy.Image(ci.context, policyContext, destRef, srcRef, &copy.Options{
		ReportWriter:   os.Stdout,
		DestinationCtx: ci.destSystemContext,
		SourceCtx:      ci.srcSystemContext,
	})
}

func targetTags(overwrite bool, src, dest []string) []string {
	if overwrite {
		return src
	}
	return missingTags(src, dest)
}

func missingTags(src, dest []string) []string {
	var result []string

	if len(dest) == 0 {
		return src
	}

	m := make(map[string]bool)
	for _, i := range dest {
		m[i] = true
	}

	for _, tag := range src {
		if !m[tag] {
			result = append(result, tag)
		}
	}
	return result
}

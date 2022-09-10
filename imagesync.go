package imagesync

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	dockerarchive "github.com/containers/image/v5/docker/archive"
	ociarchive "github.com/containers/image/v5/oci/archive"
	ocilayout "github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var Version string

var ErrInvalidTag = errors.New("invalid tag")

func Execute() error {

	app := cli.NewApp()
	app.Name = "imagesync"
	app.Usage = "Sync docker images in registries."
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "src, s",
			Usage:    "Reference for the source docker image/repository.",
			Required: true,
		},
		cli.BoolFlag{
			Name:  "src-disable-tls",
			Usage: "Disable TLS for connections to source docker registry.",
		},
		cli.StringFlag{
			Name:  "src-type",
			Usage: "[Deprecated] Type of the source docker registry",
			Value: "insecure",
		},
		cli.StringFlag{
			Name:     "dest, d",
			Usage:    "Reference for the destination docker repository.",
			Required: true,
		},
		cli.BoolFlag{
			Name:  "dest-disable-tls",
			Usage: "Disable TLS for connections to destination docker registry",
		},
		cli.StringFlag{
			Name:  "dest-type",
			Usage: "[Deprecated] Type of the destination docker registry",
			Value: "insecure",
		},
		cli.StringFlag{
			Name:  "skip-tags",
			Usage: "Comma separated list of tags to be skipped",
		},
		cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Use this to copy/override all the tags.",
		},
		cli.IntFlag{
			Name:  "max-concurrent-tags",
			Usage: "Maximum number of tags to be synced/copied in parallel.",
			Value: 1,
		},
	}

	app.Action = cli.ActionFunc(DetectAndCopyImage)

	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}

// DetectAndCopyImage will try to detect the source type and will
//
// copy the image. Detection is based on following rules if:
// - src is a directory assume it is an OCI layout.
// - src is file detect for oci-archive or docker-archive.
// - src is an image with a tag copy single image to dest.
// - none of the above then it is an entire repository sync
//	 to sync the repositories.
func DetectAndCopyImage(c *cli.Context) error {
	dest := c.String("dest")
	destRef, err := docker.ParseReference(fmt.Sprintf("//%s", dest))
	if err != nil {
		return fmt.Errorf("parsing destination ref: %w", err)
	}

	// setup copy options
	opts := copy.Options{
		ReportWriter:       os.Stdout,
		ImageListSelection: copy.CopyAllImages,
	}
	if c.String("dest-type") != "" || c.String("src-type") != "" {
		logrus.Warning("src-type/dest-type flags deprecated and will be remove in next release")
	}
	if c.String("dest-type") == "insecure" || c.Bool("dest-disable-tls") {
		opts.DestinationCtx = &types.SystemContext{DockerInsecureSkipTLSVerify: types.NewOptionalBool(true)}
	}
	if c.String("src-type") == "insecure" || c.Bool("src-disable-tls") {
		opts.SourceCtx = &types.SystemContext{DockerInsecureSkipTLSVerify: types.NewOptionalBool(true)}
	}

	ctx := context.Background()
	src := c.String("src")
	if info, err := os.Stat(src); err == nil {
		// copy oci layout
		if info.IsDir() {
			srcRef, err := ocilayout.ParseReference(src)
			if err != nil {
				return fmt.Errorf("parsing source oci ref: %w", err)
			}
			if err = copyImage(ctx, destRef, srcRef, &opts); err != nil {
				return fmt.Errorf("copy oci layout: %w", err)
			}
			logrus.Info("Image(s) sync completed.")
			return nil
		}

		// try copying oci archive with docker archive as fallback
		srcRef, _ := ociarchive.ParseReference(src)
		if err = copyImage(ctx, destRef, srcRef, &opts); err != nil {
			srcRef, err = dockerarchive.ParseReference(src)
			if err != nil {
				return fmt.Errorf("parsing source docker-archive ref: %w", err)
			}
			if err = copyImage(ctx, destRef, srcRef, &opts); err != nil {
				return fmt.Errorf("copy docker-archive layout: %w", err)
			}
		}
	} else {
		// copy single tag sync entire repository
		srcRef, err := docker.ParseReference(fmt.Sprintf("//%s", src))
		if err != nil {
			return fmt.Errorf("parsing source docker ref: %w", err)
		}
		if hasTag(src, srcRef) {
			if err = copyImage(ctx, destRef, srcRef, &opts); err != nil {
				return fmt.Errorf("copy tag: %w", err)
			}
		} else {
			if hasTag(dest, destRef) {
				return fmt.Errorf("tag shouldn't be provided in dest: %w", ErrInvalidTag)
			}
			if err = copyRepository(ctx, c, destRef, srcRef, opts); err != nil {
				return fmt.Errorf("copy repository: %w", err)
			}
		}
	}

	logrus.Info("Image(s) sync completed.")
	return nil
}

func copyRepository(ctx context.Context, cliCtx *cli.Context, destRepository, srcRepository types.ImageReference, opts copy.Options) error {
	srcTags, err := docker.GetRepositoryTags(ctx, opts.SourceCtx, srcRepository)
	if err != nil {
		return fmt.Errorf("getting source tags: %w", err)
	}

	// skip tags
	shouldSkip := cliCtx.String("skip-tags")
	if shouldSkip != "" {
		srcTags = subtract(srcTags, strings.Split(shouldSkip, ","))
	}

	var tags []string
	destTags, err := docker.GetRepositoryTags(ctx, opts.DestinationCtx, destRepository)
	if cliCtx.Bool("overwrite") || err != nil {
		tags = srcTags
	} else {
		tags = subtract(srcTags, destTags)
	}

	if len(tags) == 0 {
		logrus.Info("Image in repositories are already synced")
		os.Exit(0)
	}

	logrus.Infof("Starting image sync with total-tags=%d tags=%v source=%s destination=%s", len(tags), tags, srcRepository.DockerReference().Name(), destRepository.DockerReference().Name())

	// limit the go routines to avoid 429 on registries
	maxConcurrentTags := cliCtx.Int("max-concurrent-tags")
	numberOfConcurrentTags := maxConcurrentTags
	if len(tags) < maxConcurrentTags {
		numberOfConcurrentTags = len(tags)
	}

	// sync repository by copying each tag. Errors are ignored on purpose
	// and only warning are shown via ReportWriter for failing tags.
	var wg sync.WaitGroup
	ch := make(chan string, len(tags))
	wg.Add(numberOfConcurrentTags)
	for i := 0; i < numberOfConcurrentTags; i++ {
		go func() {
			for {
				tag, ok := <-ch
				if !ok {
					wg.Done()
					return
				}
				destTagRef, err := docker.ParseReference(fmt.Sprintf("//%s:%s", cliCtx.String("dest"), tag))
				if err != nil {
					logrus.Warnf("failed parsing dest ref: %s", err)
				}
				srcTagRef, err := docker.ParseReference(fmt.Sprintf("//%s:%s", cliCtx.String("src"), tag))
				if err != nil {
					logrus.Warnf("failed parsing src ref: %s", err)
				}
				if err = copyImage(ctx, destTagRef, srcTagRef, &opts); err != nil {
					logrus.Warnf("failed copying image: %s", err)
				}
			}
		}()
	}
	for _, tag := range tags {
		ch <- tag
	}
	close(ch)
	wg.Wait()
	return nil
}

func copyImage(ctx context.Context, destRef, srcRef types.ImageReference, opts *copy.Options) error {
	policyContext, err := signature.NewPolicyContext(&signature.Policy{
		Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()},
	})
	if err != nil {
		return fmt.Errorf("creating policy context: %w", err)
	}
	if _, err = copy.Image(ctx, policyContext, destRef, srcRef, opts); err != nil {
		return fmt.Errorf("copying image: %w", err)
	}

	return nil
}

func hasTag(ref string, imageRef types.ImageReference) bool {
	return strings.HasSuffix(imageRef.DockerReference().String(), ref)
}

func subtract(ts1 []string, ts2 []string) []string {
	var diff []string
	for _, term := range ts1 {
		if contains(ts2, term) {
			continue
		}
		diff = append(diff, term)
	}
	return diff
}

func contains(items []string, item string) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}

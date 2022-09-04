package imagesync

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	"github.com/pkg/errors"
)

type ImageCopier struct {
	src              string
	srcSystemContext *types.SystemContext

	dest              string
	destSystemContext *types.SystemContext
}

func (ic *ImageCopier) Copy(ctx context.Context, tag string) error {
	destRef, err := docker.ParseReference(fmt.Sprintf("//%s:%s", ic.dest, tag))
	if err != nil {
		return errors.WithMessagef(err, "tag=%s: parsing dest reference", tag)
	}

	srcRef, err := docker.ParseReference(fmt.Sprintf("//%s:%s", ic.src, tag))
	if err != nil {
		return errors.WithMessagef(err, "tag=%s: parsing src reference", "")
	}

	policyContext, err := signature.NewPolicyContext(&signature.Policy{
		Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()},
	})
	if err != nil {
		return errors.WithMessage(err, "creating policy context")
	}

	_, err = copy.Image(ctx, policyContext, destRef, srcRef, &copy.Options{
		ReportWriter:   os.Stdout,
		DestinationCtx: ic.destSystemContext,
		SourceCtx:      ic.srcSystemContext,
	})
	if err != nil {
		return errors.WithMessage(err, "copying image")
	}

	return nil
}

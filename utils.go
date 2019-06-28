package imagesync

import "github.com/containers/image/types"

func newSystemContextWithInsecureRegistry() (*types.SystemContext, error) {
	return &types.SystemContext{DockerInsecureSkipTLSVerify: types.NewOptionalBool(true)}, nil
}

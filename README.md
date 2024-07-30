<p>
    <a href="https://github.com/mqasimsarfraz/imagesync/actions/workflows/imagesync-ci.yaml">
        <img src="https://github.com/mqasimsarfraz/imagesync/actions/workflows/imagesync-ci.yaml/badge.svg" alt="CI"/></a>
    <a href="https://github.com/mqasimsarfraz/imagesync/actions/workflows/codeql-analysis.yml">
        <img src="https://github.com/mqasimsarfraz/imagesync/actions/workflows/codeql-analysis.yml/badge.svg" alt="codeql"/></a>
    <a href="https://mqasimsarfraz.github.io/">
        <img src="https://img.shields.io/maintenance/yes/2024.svg" alt="maintained"/></a>
    <a href="https://hub.docker.com/r/smqasims/imagesync">
        <img src="https://img.shields.io/docker/pulls/smqasims/imagesync.svg" alt="pulls"/></a>
</p>

# imagesync

A tool to copy/sync images in registries without a demon.

```bash
imagesync -h

NAME:
   imagesync - Sync images in registries.

USAGE:
   imagesync [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --src value, -s value        Reference for the source container image/repository.
   --src-strict-tls             Enable strict TLS for connections to source container registry.
   --dest value, -d value       Reference for the destination container repository.
   --dest-strict-tls            Enable strict TLS for connections to destination container registry.
   --tags-pattern value         Regex pattern to select for tag to-be synced.
   --skip-tags-pattern value    Regex pattern to exclude tags.
   --skip-tags value            Comma separated list of tags to be skipped.
   --overwrite                  Use this to copy/override all the tags.
   --max-concurrent-tags value  Maximum number of tags to be synced/copied in parallel. (default: 1)
   --help, -h                   show help
```

## Installation

### Binary

You can download the binary from [releases](https://github.com/mqasimsarfraz/imagesync/releases) page and use it directly:

```bash
VERSION=$(curl -s https://api.github.com/repos/mqasimsarfraz/imagesync/releases/latest | jq -r .tag_name)
curl -sL https://github.com/mqasimsarfraz/imagesync/releases/download/${VERSION}/imagesync_Linux_x86_64.tar.gz | sudo tar -C /usr/local/bin -xzf - imagesync
imagesync -h
```

### Docker

You can use the docker image to run `imagesync`:

```bash
VERSION=$(curl -s https://api.github.com/repos/mqasimsarfraz/imagesync/releases/latest | jq -r .tag_name)
docker run --rm -it  ghcr.io/mqasimsarfraz/imagesync:$VERSION -h
```

## Examples
Following is a list of examples with different sources. In order to try out examples with [testdata](testdata) you need to start a local [registry](https://docs.docker.com/registry/deploying/#run-a-local-registry) using:

```
docker run -d -p 5000:5000 --restart=always --name registry registry:2
```

### Docker Archive

```
imagesync  -s testdata/alpine.tar -d localhost:5000/library/alpine:3
```

### OCI Archive

```
imagesync  -s testdata/alpine-oci.tar -d localhost:5000/library/alpine:3
```

### OCI layout

```
imagesync  -s testdata/alpine-oci -d localhost:5000/library/alpine:3
```

### Image Tag

#### container image
```
imagesync  -s library/alpine:3 -d localhost:5000/library/alpine:3
```

#### helm chart
```
imagesync  -s ghcr.io/nginxinc/charts/nginx-ingress:1.3.1 -d localhost:5000/nginxinc/charts/nginx-ingress:1.3.1
```

### Entire Repository

```
imagesync  -s library/alpine -d localhost:5000/library/alpine
```

### Entire Repository (helm)

```
imagesync -s ghcr.io/nginxinc/charts/nginx-ingress -d localhost:5000/nginxinc/charts/nginx-ingress
```

## Private Registries

`imagesync` will respect the credentials stored in `~/.docker/config.json` via `docker login` etc. So in case you are
running it in a container you need to mount the path with credentials as:

```
docker run --rm -it  -v ${HOME}/.docker/config.json:/root/.docker/config.json  ghcr.io/mqasimsarfraz/imagesync:v1.2.0 -h
```

## Multi-arch images

`imagesync` supports copying multi-arch images. So in case you are copying a multi-arch image it will copy all the platforms unlike `docker pull`/`docker push` approach which only copies the platform of the host.

## Contributing/Dependencies

Following needs to be installed in order to compile the project locally:

### fedora/centos

```
dnf --enablerepo=powertools install gpgme-devel
dnf install libassuan  libassuan-devel
```

### debian/ubuntu

```
sudo apt install libgpgme-dev libassuan-dev libbtrfs-dev libdevmapper-dev pkg-config
```

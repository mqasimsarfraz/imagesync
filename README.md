<p>
    <a href="https://github.com/mqasimsarfraz/imagesync/actions/workflows/imagesync-ci.yaml">
        <img src="https://github.com/mqasimsarfraz/imagesync/actions/workflows/imagesync-ci.yaml/badge.svg" alt="CI"/></a>
    <a href="https://github.com/mqasimsarfraz/imagesync/actions/workflows/codeql-analysis.yml">
        <img src="https://github.com/mqasimsarfraz/imagesync/actions/workflows/codeql-analysis.yml/badge.svg" alt="codeql"/></a>
    <a href="https://mqasimsarfraz.github.io/">
        <img src="https://img.shields.io/maintenance/yes/2022.svg" alt="maintained"/></a>
    <a href="https://hub.docker.com/r/smqasims/imagesync">
        <img src="https://img.shields.io/docker/pulls/smqasims/imagesync.svg" alt="pulls"/></a>
</p>

# imagesync

A tool to copy/sync container images in registries without a demon.

## Command

```
docker run --rm -it smqasims/imagesync:v1.1.0 -h
```

or 

```
imagesync -h
```

## Usage

```
NAME:
   imagesync - Sync container images in registries.

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
   --skip-tags value            Comma separated list of tags to be skipped.
   --overwrite                  Use this to copy/override all the tags.
   --max-concurrent-tags value  Maximum number of tags to be synced/copied in parallel. (default: 1)
   --help, -h                   show help
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

```
imagesync  -s library/alpine:3 -d localhost:5000/library/alpine:3
```

### Entire Repository

```
imagesync  -s library/alpine -d localhost:5000/library/alpine
```

## Private Registries

`imagesync` will respect the credentials stored in `~/.docker/config.json` via `docker login` etc. So in case you are
running it in a container you need to mount the path with credentials as:

```
docker run --rm -it  -v ${HOME}/.docker/config.json:/root/.docker/config.json  smqasims/imagesync:v1.1.0 -h
```

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

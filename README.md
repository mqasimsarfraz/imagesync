<p>
    <a href="https://hub.docker.com/r/smqasims/imagesync" alt="Pulls">
        <img src="https://img.shields.io/docker/pulls/smqasims/imagesync.svg" /></a>
    <a href="https://mqasimsarfraz.github.io/" alt="Maintained">
        <img src="https://img.shields.io/maintenance/yes/2022.svg" /></a>

</p>

# imagesync

A tool to copy/sync docker images in registries without docker demon.

## Command

```
docker run --rm -it smqasims/imagesync:v1.0.2 -h
```

## Usage

```
NAME:
   imagesync - Sync docker images between repositories.

USAGE:
   imagesync [global options] command [command options] [arguments...]

VERSION:
   1.0.2

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --src value, -s value   Reference for the source docker registry.
   --src-type value        Type of the source docker registry (default: "insecure")
   --dest value, -d value  Reference for the destination docker registry.
   --dest-type value       Type of the destination docker registry (default: "insecure")
   --overwrite             Use this to copy/override all the tags.
   --help, -h              show help
   --version, -v           print the version

```

## Private Registries

`imagesync` will respect the credentials stored in `~/.docker/config.json` via `docker login` etc. So in case you are
running it in a container you need to mount the path with credentials as:

```
docker run --rm -it  -v ${HOME}/.docker/config.json:/root/.docker/config.json  smqasims/imagesync:v1.0.2 -h
```

## Contributing/Dependencies

Following needs to be installed in order to compile the project locally:

### fedora/centos

```
dnf --enablerepo=powertools install gpgme-devel
dnf install libassuan  libassuan-devel
```

### debain/ubuntu

```
sudo apt install libgpgme-dev libassuan-dev libbtrfs-dev libdevmapper-dev pkg-config
```

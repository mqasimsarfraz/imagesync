<p align="left">
    <a href="https://hub.docker.com/r/smqasims/imagesync/builds" alt="Build">
        <img src="https://img.shields.io/docker/cloud/build/smqasims/imagesync.svg" /></a>
    <a href="https://hub.docker.com/r/smqasims/imagesync" alt="Pulls">
        <img src="https://img.shields.io/docker/pulls/smqasims/imagesync.svg" /></a>
    <a href="https://mqasimsarfraz.github.io/" alt="Maintained">
        <img src="https://img.shields.io/maintenance/yes/2020.svg" /></a>
        
</p>

# imagesync
A tool to copy/sync docker images between registries without docker deamon.
## Command
```
docker run --rm -it smqasims/imagesync -h
```
## Usage
```
NAME:
   imagesync - Sync docker images between repositories.

USAGE:
   imagesync [global options] command [command options] [arguments...]

VERSION:
   0.0.0

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
## Note
Currently it isn't possible to specify credentials for a private registry. But you can still use the `imagesync` if you are already logged in using docker cli. It will respect the credentials stored `~/.docker/config.json` via `docker login`. So in case you are running it in a container you need to mount this path:
```
docker run --rm -it  -v ${HOME}/.docker/config.json:/root/.docker/config.json  smqasims/imagesync -h
```
## Todo(s):
- Add support to choose image policies.

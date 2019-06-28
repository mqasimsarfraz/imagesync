<p align="left">
    <a href="https://hub.docker.com/r/smqasims/imagesync/builds" alt="Build">
        <img src="https://img.shields.io/docker/cloud/build/smqasims/imagesync.svg" /></a>
    <a href="https://hub.docker.com/r/smqasims/imagesync" alt="Pulls">
        <img src="https://img.shields.io/docker/pulls/smqasims/imagesync.svg" /></a>
    <a href="https://mqasimsarfraz.github.io/" alt="Maintained">
        <img src="https://img.shields.io/maintenance/yes/2019.svg" /></a>
        
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
## Example:
```

```
## Todo(s):
- Add support for further registries e.g AWS ECR, Jfrog etc
- Add support to choose image policies.

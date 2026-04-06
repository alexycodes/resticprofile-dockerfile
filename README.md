# alexycodes/resticprofile

[![Docker Pulls](https://img.shields.io/docker/pulls/alexycodes/resticprofile?label=Docker%20Pulls&logo=docker)](https://hub.docker.com/r/alexycodes/resticprofile)

A [resticprofile](https://github.com/creativeprojects/resticprofile) image built for long-lived containers with database and backup utilities.

It overrides the official image's entrypoint and command to execute `crond` on start, so containers run indefinitely and can be used to schedule tasks.

Images are available on Docker Hub at [alexycodes/resticprofile](https://hub.docker.com/r/alexycodes/resticprofile) and can be installed easily on Unraid from the app store.

## Usage

### Docker

```sh
docker run -d \
  --name resticprofile \
  -e TZ=Your/Timezone \
  -v /path/to/profiles:/resticprofile \
  -v /path/to/passfiles:/pass \
  -v /path/to/data:/data \
  -v /path/to/crontab:/crontab \ # optional: mount a crontab to schedule tasks
  alexycodes/resticprofile
```

### Docker Compose

```yaml
services:
  resticprofile:
    image: alexycodes/resticprofile
    environment:
      - TZ: Your/Timezone
    volumes:
      - /path/to/profiles:/resticprofile
      - /path/to/passfiles:/pass
      - /path/to/data:/data
      - /path/to/crontab:/crontab # optional: mount a crontab to schedule tasks
```

### Run a one-off command

The following example overrides the default `crond` command to run a resticprofile command directly and exit.

```sh
docker run --rm \
  -e TZ=Your/Timezone \
  -v /path/to/profiles:/resticprofile \
  -v /path/to/passfiles:/pass \
  -v /path/to/data:/data \
  alexycodes/resticprofile \
  resticprofile backup
```

## Configuration

### profiles.yaml

The following is a partial example demonstrating how to reference paths inside the container. See the [full documentation](https://creativeprojects.github.io/resticprofile/configuration/index.html) for all configuration options.

```yaml
# /resticprofile/profiles.yaml

default:
  password-file: /pass/example

  backup:
    source:
      - /data
```

### crontab

Optionally, a [crontab](https://linuxhandbook.com/crontab/) configuration file can be mounted at `/crontab` to schedule tasks. The entrypoint automatically loads it on start.

The following example runs the `default` profile's backup command every day at midnight.

```
0 0 * * * cd /resticprofile && resticprofile backup
```

## Installed packages

Images have the following packages installed on top of the official resticprofile images:

- `acl`
- `bash`
- `btrfs-progs`
- `ca-certificates`
- `curl`
- `findmnt`
- `fuse`
- `libxxhash`
- `logrotate`
- `lz4-libs`
- `mariadb-client`
- `mariadb-connector-c`
- `mongodb-tools`
- `openssl`
- `postgresql-client`
- `sqlite`
- `sshfs`
- `tzdata`
- `xxhash`

# ilium007/resticprofile

![Static Badge](https://img.shields.io/badge/github-repo-blue?logo=github&label=github&link=https%3A%2F%2Fghcr.io%2Filium007%2Fresticprofile-dockerfile)

A [resticprofile](https://github.com/creativeprojects/resticprofile) image built for long-lived containers with database and backup utilities. Forked from the excellent work at [alexycodes/resticprofile-dockerfile](https://github.com/alexycodes/resticprofile-dockerfile).

It overrides the official image's entrypoint and command to execute `crond` on start, so containers run indefinitely and can be used to schedule tasks.

## Scheduling tasks

A [crontab](https://linuxhandbook.com/crontab/) configuration file can be mounted inside the container at `/crontab` to automatically schedule tasks.

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
- `rclone`
- `sqlite`
- `sshfs`
- `tzdata`
- `xxhash`

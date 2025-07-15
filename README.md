# alexycodes/resticprofile

[![Docker Pulls](https://img.shields.io/docker/pulls/alexycodes/resticprofile?label=Docker%20Pulls&logo=docker)](https://hub.docker.com/r/alexycodes/resticprofile)

A [resticprofile](https://github.com/creativeprojects/resticprofile) image built for long-lived containers with database and backup utilities.

It overrides the official image's entrypoint and command to execute `crond` on start, so containers run indefinitely and can be used to schedule tasks.

## Scheduling tasks

A [crontab](https://linuxhandbook.com/crontab/) configuration file can be mounted inside the container at `/crontab` to automatically schedule tasks.

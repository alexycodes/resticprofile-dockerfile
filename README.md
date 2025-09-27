# alexycodes/resticprofile

A resticprofile image with its entrypoint and command overridden to execute `crond` on start. Containers therefore run indefinitely and can be used to schedule tasks.

## Scheduling tasks

A [crontab](https://linuxhandbook.com/crontab/) configuration file can be mounted inside the container at `/crontab` to automatically schedule tasks.

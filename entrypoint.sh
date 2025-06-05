#!/bin/sh

# If a user-provided /crontab exists, load it
if [ -f /crontab ]; then
  crontab /crontab
fi

exec "$@"

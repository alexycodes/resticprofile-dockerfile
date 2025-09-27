FROM creativeprojects/resticprofile:0.12.0

RUN apk add --update --no-cache \
        acl                     \
        bash                    \
        btrfs-progs             \
        ca-certificates         \
        curl                    \
        findmnt                 \
        fuse                    \
        libxxhash               \
        logrotate               \
        lz4-libs                \
        mariadb-client          \
        mariadb-connector-c     \
        mongodb-tools           \
        openssl                 \
        postgresql-client       \
        sqlite                  \
        sshfs                   \
        tzdata                  \
        xxhash

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["crond", "-f"]

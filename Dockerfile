FROM creativeprojects/resticprofile:0.31.0

COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

CMD ["crond", "-f"]

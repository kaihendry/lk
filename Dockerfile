FROM scratch
COPY lk /
COPY i /srv
EXPOSE 3000
ENTRYPOINT ["/lk", "--port", "3000", "/srv"]

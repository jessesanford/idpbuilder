FROM alpine:latest
RUN invalid-command-that-does-not-exist
INVALID SYNTAX HERE
CMD ["echo", "this will fail"]
FROM ubuntu:latest
WORKDIR /usr/src/app

COPY . .
COPY pawang-be /usr/local/bin/pawang-be

EXPOSE 1234
CMD ["pawang-be"]
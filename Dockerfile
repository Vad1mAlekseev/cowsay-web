FROM golang:1.16

WORKDIR /app
COPY . .

ENV PATH="/usr/games:${PATH}"

RUN apt update && \
    apt install cowsay fortune -y && \
    make

CMD ["./build/cowsay-web"]

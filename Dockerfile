FROM ubuntu

WORKDIR /usr/src/app
COPY . .

RUN apt-get update && \
    apt-get install -y curl && \
    curl -sL https://deb.nodesource.com/setup_12.x | bash - && \
    apt-get install -y nodejs cowsay fortune && \
    npm install --only=prod

# "cowsay" installs to /usr/games
ENV PATH $PATH:/usr/games

CMD ["node", "main.js"]

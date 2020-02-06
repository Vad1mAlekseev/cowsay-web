FROM ubuntu

RUN apt-get update && \
    apt-get install -y curl && \
    curl -sL https://deb.nodesource.com/setup_12.x | bash - && \
    apt-get install -y nodejs cowsay fortune

# "cowsay" installs to /usr/games
ENV PATH $PATH:/usr/games

WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .

CMD ["node", "main.js"]

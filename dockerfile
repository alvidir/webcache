FROM node:latest

LABEL maintainer="Hector Morales <hector.morales.carnice@gmail.com>"
LABEL repo-url="https://github.com/alvidir/unsplash"
LABEL version="alpha"

# Create app directory
WORKDIR /usr/src/app

# installing dependencies
COPY package*.json ./
RUN npm install
# RUN npm ci --only=production

# get source code
COPY . .

CMD [ "node", "server.js" ]
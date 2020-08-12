FROM node:latest

LABEL maintainer="Hector Morales <hector.morales.carnice@gmail.com>"
LABEL repo-url="https://github.com/alvidir/unsplash-api"
LABEL version="alpha"

# Create app directory
WORKDIR /usr/src/app

# installing dependencies
COPY package.json .
RUN npm install
# RUN npm ci --only=production

# get source code
COPY . .

# typescript to js
# RUN npm run tsc

# Execute server from typescript
CMD ["npx", "ts-node", "src/app.ts"]
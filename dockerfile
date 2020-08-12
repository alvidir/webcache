FROM node:lts-alpine as builder

LABEL maintainer="Hector Morales <hector.morales.carnice@gmail.com>"
LABEL repo-url="https://github.com/alvidir/unsplash-api"
LABEL version="alpha"

WORKDIR /srv

# installing dependencies
COPY package*.json /srv/
RUN npm ci

COPY tsconfig.json /srv/
COPY src /srv/src/
RUN npm run tsc

RUN npm ci --production

# Execute server from typescript
# CMD ["npx", "ts-node", "src/app.ts"]

# start new image from scratch
FROM alpine:latest

RUN apk add nodejs --no-cache

WORKDIR /srv
COPY --from=builder /srv/node_modules /srv/node_modules
COPY --from=builder /srv/build /srv/

CMD ["node", "app.js"]
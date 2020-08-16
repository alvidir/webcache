FROM node:latest as builder

LABEL maintainer="Hector Morales <hector.morales.carnice@gmail.com>"
LABEL repo-url="https://github.com/alvidir/unsplash-api"
LABEL version="alpha"

WORKDIR /app

# installing dependencies
COPY package*.json ./
RUN npm install

COPY . .
# Index.ts may does nothing by itself
# However its presence in the project's root ensures that the build script
# will keep the directory structure of this project.
RUN touch index.ts
RUN npm run build

# Execute server from typescript
# CMD ["npx", "ts-node", "src/app.ts"]

# start new image from scratch
FROM node:lts-alpine

RUN apk add nodejs --no-cache

WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/build/* ./
COPY --from=builder /app/.env ./

CMD ["node", "app.js"]
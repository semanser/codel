# STEP 1: Build the frontend
FROM node:21-slim as fe-build

ENV NODE_ENV=production
ENV VITE_API_URL=localhost:3000

WORKDIR /frontend

COPY ./backend/graph/schema.graphqls ../backend/graph/

COPY frontend/ .

# --production=false is required because we want to install the @graphql-codegen/cli package (and it's in the devDependencies)
# https://classic.yarnpkg.com/lang/en/docs/cli/install/#toc-yarn-install-production-true-false
RUN yarn install --frozen-lockfile --production=false
RUN ls -la /frontend
RUN yarn build

# STEP 2: Build the backend
FROM golang:1.22-alpine as be-build
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

WORKDIR /backend

COPY backend/ .

RUN go mod download

RUN go build -ldflags='-extldflags "-static"' -o /app

# STEP 3: Build the final image
FROM alpine:3.14

COPY --from=be-build /app /app
COPY --from=fe-build /frontend/dist /fe

# Install sqlite3

RUN apk add --no-cache sqlite

CMD /app

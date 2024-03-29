FROM node:21-slim

ENV NODE_ENV=production

WORKDIR /frontend

COPY ./backend/graph/schema.graphqls ../backend/graph/

COPY frontend/ .

# --production=false is required because we want to install the @graphql-codegen/cli package (and it's in the devDependencies)
# https://classic.yarnpkg.com/lang/en/docs/cli/install/#toc-yarn-install-production-true-false
RUN yarn install --frozen-lockfile --production=false
RUN ls -la /frontend
RUN yarn build

FROM node:gallium-bullseye-slim AS base

ENV NODE_ENV production

RUN mkdir -p /root/.docker/
RUN chmod +x /root/.docker

WORKDIR /instant

# install curl
RUN apt-get update; apt-get install -y curl

# install docker engine
RUN curl -sSL https://get.docker.com/ | sh

# install docker-compose binary
RUN curl -L "https://github.com/docker/compose/releases/download/v2.16.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose

# remove orphan container warning
ENV COMPOSE_IGNORE_ORPHANS=1

# default ENV vars for instant container
ENV CLUSTERED_MODE false

# install node deps
ADD package.json .
ADD yarn.lock .
RUN yarn --production --frozen-lockfile

# add entrypoint script
ADD instant.ts .

# add util function scripts
ADD utils ./utils

# add schema
ADD schema ./schema

ENTRYPOINT [ "yarn", "instant" ]

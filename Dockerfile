FROM node:gallium-bullseye-slim AS base

WORKDIR /instant

# install curl
RUN apt-get update; apt-get install -y curl

# install docker engine
RUN curl -sSL https://get.docker.com/ | sh

# install docker-compose binary
RUN curl -L "https://github.com/docker/compose/releases/download/1.25.5/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
RUN chmod +x /usr/local/bin/docker-compose

# remove orphan container warning
ENV COMPOSE_IGNORE_ORPHANS=1

# install node deps
ADD package.json .
ADD yarn.lock .
RUN yarn --prod

# add entrypoint script
ADD instant.ts .

ENTRYPOINT [ "yarn", "instant" ]
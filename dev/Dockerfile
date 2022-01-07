ARG GO_VERSION
FROM golang:$GO_VERSION

ARG PG_MAJOR
ENV APP_PATH=/fylerx

# Install dependencies
COPY Aptfile /tmp/Aptfile
RUN apt-get update -qq && DEBIAN_FRONTEND=noninteractive apt-get -yq dist-upgrade && \
  DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends \
  build-essential \
  postgresql-client-$PG_MAJOR \
  $(grep -Ev '^\s*#' /tmp/Aptfile | xargs) && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
  truncate -s 0 /var/log/*log

# Create a directory for the app code
RUN mkdir -p ${APP_PATH}

WORKDIR ${APP_PATH}

FROM golang:1.16.3-alpine

# Add make command
RUN apk add --no-cache make bash

WORKDIR /app
COPY . /app
RUN make

WORKDIR /appbin
RUN cp /app/github-contributors-action /appbin/
RUN rm -r /app

ENV PATH=${PATH}:/appbin

CMD ["github-contributors-action"]

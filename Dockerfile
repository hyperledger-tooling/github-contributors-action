FROM golang:1.16.3
WORKDIR /app
COPY . /app
RUN go build
RUN rm *.go go.*
ENV PATH=${PATH}:/app

CMD ["github-contributors-action"]
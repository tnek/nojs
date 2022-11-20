FROM debian:bookworm as builder
MAINTAINER tnek

RUN apt-get update && apt-get -y install ca-certificates git golang
ADD keys/deploy_id_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
COPY . .

RUN go mod tidy
RUN go mod vendor
RUN go mod download

RUN go build -installsuffix 'static' -o /notes-site .

FROM debian:bookworm AS final
COPY --from=builder /notes-site ./notes-site
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/assets ./assets

RUN apt-get update && apt-get install -y curl xvfb openjdk-11-jre software-properties-common unzip wget libc6-amd64-cross libc6-dev firefox-esr 

# Geckodriver
RUN VERSION=$(curl -sL https://api.github.com/repos/mozilla/geckodriver/releases/latest | grep tag_name | cut -d '"' -f 4) && curl -sL "https://github.com/mozilla/geckodriver/releases/download/$VERSION/geckodriver-$VERSION-linux-aarch64.tar.gz" | tar -xz -C /usr/local/bin

# Get selenium's jar
ENV SELENIUM_JAR_ADDR=https://github.com/SeleniumHQ/selenium/releases/download/selenium-3.141.59/selenium-server-standalone-3.141.59.jar
RUN curl -sL $SELENIUM_JAR_ADDR > /usr/local/bin/selenium-server.jar

CMD ["./notes-site"]

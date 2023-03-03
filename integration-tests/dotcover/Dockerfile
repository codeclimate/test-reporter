FROM mcr.microsoft.com/dotnet/sdk:7.0

# Install GoLang
RUN curl -O https://dl.google.com/go/go1.15.linux-amd64.tar.gz
RUN tar -xzf go1.15.linux-amd64.tar.gz
RUN mv go /usr/local

ENV PATH $PATH:/usr/local/go/bin
ENV GOBIN="/usr/local/go/bin"
RUN go version

ENV GOPATH /go
RUN mkdir $GOPATH
ENV PATH $PATH:/go/bin

ENV CCTR=$GOPATH/src/github.com/codeclimate/test-reporter
RUN mkdir -p $CCTR
WORKDIR $CCTR
COPY . .
RUN go install -v

ENV PATH $PATH:/root/.dotnet/tools
RUN dotnet tool install JetBrains.dotCover.GlobalTool -g --version "2022.3.2"

# Clone .NET example repo and run test
RUN git clone https://github.com/codeclimate/dot-net-coverage-test.git
WORKDIR dot-net-coverage-test
RUN dotnet build
RUN dotnet dotcover test --dcReportType=DetailedXML --dcOutput="dotcover.xml" --no-build

RUN echo "testing" > ignore.me && \
  git config --global user.email "you@example.com" && \
  git config --global user.name "Your Name" && \
  git add ignore.me && \
  git commit -m "testing"

ENV CC_TEST_REPORTER_ID=49a59a1849364524250d544e098b5d987335376cdb739ea7649c9f8bce968e3b
RUN test-reporter format-coverage -d -t dotcover
RUN cat coverage/codeclimate.json
RUN test-reporter upload-coverage -d -s 2

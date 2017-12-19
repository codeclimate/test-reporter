FROM elixir:1.4.2

RUN curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
RUN tar -xvf go1.8.linux-amd64.tar.gz
RUN mv go /usr/local

ENV PATH $PATH:/usr/local/go/bin
RUN go version

ENV GOPATH /go
RUN mkdir $GOPATH
ENV PATH $PATH:/go/bin

ENV CCTR=$GOPATH/src/github.com/codeclimate/test-reporter
RUN mkdir -p $CCTR
WORKDIR $CCTR
COPY . .
RUN go install -v

RUN git clone https://github.com/ale7714/excoveralls.git
WORKDIR excoveralls

RUN echo "testing" > ignore.me
RUN git config --global user.email "you@example.com"
RUN git config --global user.name "Your Name"
RUN git add ignore.me
RUN git commit -m "testing"
ENV MIX_ENV=test
RUN mix do local.hex --force, local.rebar --force, deps.get, compile, coveralls.json


ENV CC_TEST_REPORTER_ID=871a78f3199b36b65af5e00b9395319d971d3c85e467450b53b386335a60a261

RUN test-reporter after-build -s 2 -d

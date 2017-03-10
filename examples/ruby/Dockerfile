FROM ruby:latest

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

RUN git clone https://github.com/codeclimate/ruby-test-reporter.git
WORKDIR ruby-test-reporter
RUN bundle
RUN bundle exec rake

RUN REPO_TOKEN=c4881e09870b0fac1291c93339b36ffe36210a2645c1ad25e52d8fda3943fb4d test-reporter format

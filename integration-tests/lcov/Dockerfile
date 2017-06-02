FROM node:latest

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

RUN git clone https://github.com/codeclimate/javascript-test-reporter.git
WORKDIR javascript-test-reporter
RUN npm install

RUN echo "testing" > ignore.me
RUN git config --global user.email "you@example.com"
RUN git config --global user.name "Your Name"
RUN git add ignore.me
RUN git commit -m "testing"
RUN npm test; exit 0

ENV CC_TEST_REPORTER_ID=c4881e09870b0fac1291c93339b36ffe36210a2645c1ad25e52d8fda3943fb4d
RUN test-reporter format-coverage -d
RUN cat coverage/codeclimate.json
RUN test-reporter upload-coverage -d -s 2

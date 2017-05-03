FROM connectmedica/php7.1-xdebug:latest

RUN apt-get update
RUN apt-get install -y git wget

RUN wget https://phar.phpunit.de/phpunit-6.1.phar
RUN chmod +x phpunit-6.1.phar
RUN mv phpunit-6.1.phar /usr/local/bin/phpunit
RUN phpunit --version

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

RUN git clone https://github.com/codeclimate/php-test-reporter.git
WORKDIR php-test-reporter

RUN echo "testing" > ignore.me
RUN git config --global user.email "you@example.com"
RUN git config --global user.name "Your Name"
RUN git add ignore.me
RUN git commit -m "testing"
RUN composer install
RUN vendor/bin/phpunit --coverage-clover=build/logs/clover.xml

ENV CC_TEST_REPORTER_ID=c4881e09870b0fac1291c93339b36ffe36210a2645c1ad25e52d8fda3943fb4d
RUN test-reporter format-coverage -d
RUN test-reporter upload-coverage -d -s 2

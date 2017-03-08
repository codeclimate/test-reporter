.PHONY: test

IMAGE_NAME ?= codeclimate/test-reporter-acceptance

repos:
	mkdir -p ./repos
	git clone https://github.com/codeclimate/ruby-test-reporter ./repos/ruby-test-reporter
	@echo "now do JavaScript"
	@echo "now do Python"
	@echo "now do PHP"


files/bin/cc-test-reporter:
	@echo "Build ../ and put the executable in ./files/bin"
	touch files/bin/cc-test-reporter

image: files/bin/cc-test-reporter
	docker build --tag "$(IMAGE_NAME)" .

test: image
	docker run --rm --volume "$(PWD)"/repos:/opt/cc/repos "$(IMAGE_NAME)"

all: repos test

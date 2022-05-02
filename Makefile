# Copyright 2021 Hyperledger Community
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := /bin/bash
version=$(shell cat VERSION)
LDFLAGS=-ldflags "-X main.AppVersion=$(version)"
format_output=$(shell gofmt -l .)

.PHONY: all
all: clean build

clean:
	rm -f github-contributors-action

build: lint-check unit-test
	go build -o github-contributors-action $(LDFLAGS) ./cmd

unit-test:
	CGO_ENABLED=0 go test -v ./...

lint-check:
	@[ "$(format_output)" == "" ] || exit -1

format:
	go fmt ./...

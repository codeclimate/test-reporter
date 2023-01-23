# Copyright 2017 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

VERSION=v1

podman-prod: Dockerfile
	podman build -f Dockerfile --tag=gcr.io/symbolic-datum-552/tip:$(VERSION) .
podman-dev: Dockerfile
	podman build -f Dockerfile --tag=gcr.io/go-dashboard-dev/tip:$(VERSION) .

push-prod: podman-prod
	gcloud podman push -- gcr.io/symbolic-datum-552/tip:$(VERSION)
push-dev: podman-dev
	gcloud podman push -- gcr.io/go-dashboard-dev/tip:$(VERSION)

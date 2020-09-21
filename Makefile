# +-------------------------------------------------------------------------
# | Copyright (C) 2018 Yunify, Inc.
# +-------------------------------------------------------------------------
# | Licensed under the Apache License, Version 2.0 (the "License");
# | you may not use this work except in compliance with the License.
# | You may obtain a copy of the License in the LICENSE file, or at:
# |
# | http://www.apache.org/licenses/LICENSE-2.0
# |
# | Unless required by applicable law or agreed to in writing, software
# | distributed under the License is distributed on an "AS IS" BASIS,
# | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# | See the License for the specific language governing permissions and
# | limitations under the License.
# +-------------------------------------------------------------------------

.PHONY: all disk

IMAGE=csiplugin/csi-neonsan
TAG=canary
IMAGE_UBUNTU=csiplugin/csi-neonsan-ubuntu
TAG_UBUNTU=canary
IMAGE_CENTOS=csiplugin/csi-neonsan-centos
TAG_CENTOS=canary
ROOT_PATH=$(pwd)
PACKAGE_LIST=./cmd/... ./pkg/...

container:
	docker build -t ${IMAGE}:${TAG} -f deploy/neonsan/docker/Dockerfile  .

container-ubuntu:
	docker build -t ${IMAGE_UBUNTU}:${TAG_UBUNTU} -f deploy/neonsan/docker/ubuntu/Dockerfile  .

container-centos:
	docker build -t ${IMAGE_CENTOS}:${TAG_CENTOS} -f deploy/neonsan/docker/centos/Dockerfile  .

mod:
	go build ./...
	go mod download
	go mod tidy
	go mod vendor

test:
	go test -cover -mod=vendor -gcflags=-l -count=1 ./pkg/...

clean:
	go clean -r -x ./...
	rm -rf ./_output

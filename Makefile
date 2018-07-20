PROJ = mam-go-lib
docker_registry = artifactory.viaplay.cloud:5000

IMAGE_TAG ?= localdev
NOW = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_TAG = $(shell git describe --tags)
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_COMMENT = $(shell git log -1 --pretty=%B | head -1)
IMAGE_NAME = $(docker_registry)/$(PROJ):$(IMAGE_TAG)

VER = ./buildtarget/build.version
LDFLAGS = -extldflags -static -s -w

build:
	echo "now = ${NOW}" > ${VER}
	echo "project = ${PROJ}" >> ${VER}
	echo "imageName = ${IMAGE_NAME}" >> ${VER}
	echo "imageTag = ${IMAGE_TAG}" >> ${VER}
	echo "gitBranch = ${GIT_BRANCH}" >> ${VER}
	echo "gitTag = ${GIT_TAG}" >> ${VER}
	echo "gitCommit = ${GIT_COMMIT}" >> ${VER}
	echo "gitComment = ${GIT_COMMENT}" >> ${VER}
	@echo --- build version ---
	cat ${VER}
	@echo ---------------------
	env CGO_ENABLED=0 vgo build -ldflags "${LDFLAGS}" -o ./buildtarget/mam-golib-example ./cmd/mam-golib-example

clean:
	rm -rf ./buildtarget/*

docker-build:
	docker build --build-arg IMAGE_TAG=${IMAGE_TAG} -t ${IMAGE_NAME} .

docker-push: docker-build
	docker push ${IMAGE_NAME}

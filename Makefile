BALLOT_IMG=ballot
COMMITID := $(shell git rev-parse HEAD)
ifndef IMAGE_TAG
  IMAGE_TAG=latest
endif
CLUSTER_IP := $(shell ping -W2 -n -q -c1 current-cluster-roost.io  2> /dev/null | awk -F '[()]' '/PING/ { print $$2}')

# HOSTNAME := $(shell hostname)
.PHONY: all dockerise helm-deploy
all: dockerise helm-deploy

.PHONY: test
test: unittest test-ballot

.PHONY: test-ballot
test-ballot:
	@echo "========================="
	@echo "Running Artillery Test"
	@echo "========================="
	docker run --network="host" --rm -it -v ${PWD}/ballot/test:/scripts \
   zbio/artillery-custom \
   run -e unit /scripts/test.yaml

.PHONY: pre-dockerise
pre-dockerise:
	docker pull golang:1.19.3-alpine3.16
	docker pull alpine:3.16
	docker pull node:14.21.1-alpine3.16
	docker pull nginx:stable-alpine

.PHONY: dockerise
dockerise: pre-dockerise build-ballot

.PHONY: build-ballot
build-ballot:
ifdef DOCKER_HOST
	docker -H ${DOCKER_HOST} build -t ${BALLOT_IMG}:${COMMITID} -f ballot/Dockerfile ballot
	docker -H ${DOCKER_HOST} tag ${BALLOT_IMG}:${COMMITID} ${BALLOT_IMG}:${IMAGE_TAG}
else
	docker build -t ${BALLOT_IMG}:${COMMITID} -f ballot/Dockerfile ballot
	docker tag ${BALLOT_IMG}:${COMMITID} ${BALLOT_IMG}:${IMAGE_TAG}
endif

.PHONY: push
push:
	docker tag ${BALLOT_IMG}:${IMAGE_TAG} zbio/${BALLOT_IMG}:${IMAGE_TAG}
	docker push zbio/${BALLOT_IMG}:${IMAGE_TAG}

.PHONY: deploy
deploy:
	kubectl apply -f ballot/ballot.yaml
	
.PHONY: helm-deploy
helm-deploy: 
ifeq ($(strip $(CLUSTER_IP)),)
	@echo "UNKNOWN_CLUSTER_IP: failed to resolve current-cluster-roost.io to an valid IP"
	@exit 1;
endif
		helm install vote helm-vote --set clusterIP=$(CLUSTER_IP)
		
.PHONY: helm-undeploy
helm-undeploy:
		-helm uninstall vote

.PHONY: clean
clean: helm-undeploy

.PHONY: unittest
unittest:
	@echo "==================="
	@echo "Running Unit Test"
	@echo "==================="
	docker build --target ballottest -t ${BALLOT_IMG}:${IMAGE_TAG}-uniitest -f ballot/Dockerfile ballot
	docker run --rm ${BALLOT_IMG}:${IMAGE_TAG}-uniitest ballot.test -test.v -test.parallel 2 -test.count 50 -test.failfast
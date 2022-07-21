BALLOT_IMG=ballot
COMMITID := $(shell git rev-parse HEAD)
ifndef IMAGE_TAG
  IMAGE_TAG=latest
endif
CLUSTER_IP := $(shell ping -W2 -n -q -c1 current-cluster-roost.io  2> /dev/null | awk -F '[()]' '/PING/ { print $$2}')

# HOSTNAME := $(shell hostname)
.PHONY: all
all: dockerise helm-deploy

.PHONY: test
test: test-ballot

.PHONY: test-ballot
test-ballot:
	echo "Test Ballot"
	docker run --network="host" --rm -it -v ${PWD}/ballot/test:/scripts \
   zbio/artillery-custom \
   run -e unit /scripts/test.yaml

.PHONY: dockerise
dockerise: build-ballot

.PHONY: build-ballot
build-ballot:
ifdef DOCKER_HOST
	docker -H ${DOCKER_HOST} build -t ${BALLOT_IMG}:${COMMITID} -f ballot/Dockerfile ballot
	docker -H ${DOCKER_HOST} tag ${BALLOT_IMG}:${COMMITID} ${BALLOT_IMG}:${IMAGE_TAG}
else
	docker build -t ${BALLOT_IMG}:${IMAGE_TAG} -f ballot/Dockerfile ballot
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

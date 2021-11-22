attack:
	vegeta attack -cert ./pki/admin.crt -key ./pki/admin.key -duration=1s -root-certs=./pki/ca.crt


.dockerfilestamp: Dockerfile
	docker build . -t vegeta


SUITE ?= tunneled_mode.conf
DURATION ?= 120s
EXTRA_ARGS ?= ""

suite: .dockerfilestamp
	docker run -e DURATION=$(DURATION)  --rm -v $(PWD)/kubeconfigs/$(SUITE):/kubeconfig -v $(PWD)/reports/report_$(SUITE):/report vegeta

.PHONY: all
all: .dockerfilestamp
	 DURATION=$(DURATION) ./run_all.sh

clean:
	rm -rf reports
	rm -rf .dockerfilestamp
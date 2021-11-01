#!/bin/bash
set -ex

CONFIGS=`ls kubeconfigs`
mkdir -p reports
for SUITE in $CONFIGS
do 
    docker run -e DURATION=${DURATION} --rm -v `pwd`/kubeconfigs/${SUITE}:/kubeconfig -v `PWD`/reports/report_${SUITE}:/report vegeta
done
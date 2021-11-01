#!/bin/bash
set -ex

/bin/create_targets -kubeconfig=$1 > /tmp/output

DIR=`cat /tmp/output|tail -n 1`

vegeta \
    attack \
    -duration=${DURATION} \
    -cert ./${DIR}/cert.cert \
    -key ./${DIR}/cert.key \
    -root-certs=./${DIR}/ca.cert \
    -targets ./${DIR}/targets > /report/raw

vegeta plot /report/raw > /report/index.html
vegeta report /report/raw > /report/report
vegeta report --type=json /report/raw > /report/report.json
cp -r ${DIR} /report/suite
cp /tmp/output /report/suite/gen.out
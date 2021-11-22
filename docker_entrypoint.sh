#!/bin/bash
set -ex

/bin/create_targets -kubeconfig=$1 > /tmp/output

DIR=`cat /tmp/output|tail -n 1`


cat ${DIR}/targets
vegeta \
    attack -keepalive=false \
    -duration=${DURATION} \
    -cert ${DIR}/cert.cert \
    -key ${DIR}/cert.key \
    -root-certs=${DIR}/ca.cert \
    -targets ${DIR}/targets > ${DIR}/raw

vegeta plot ${DIR}/raw > ${DIR}/index.html
vegeta report ${DIR}/raw > ${DIR}/report
vegeta report --type=json ${DIR}/raw > ${DIR}/report.json
cp /tmp/output ${DIR}/create_targets.out

mv ${DIR} /report/${DIR}
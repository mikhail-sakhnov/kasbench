# Benchmarks for kube-api
Package has tool to generate workload for vegeta based kube-api stress tests.
## Running

### Run single
Put the kubeconfig file under the kubeconfigs dir and pass it to makefile:
```bash
$ k0s kubeconfig admin > kubeconfigs/testsuite.conf
$ SUITE=testsuite.conf make suite
$ # results available under the current working dir
$ ls reports/report_testsuite.conf
```

### Run all
```bash
$ make all
$ ls reports
```

## Reading th report

Each report includes:
- `index.html` with interactive chart
- `raw` binary data for vegeta
- `report` default vegeta text report
- `report.json` same data in json

Directory `suite` includes the certificates for the cluster and instructions how to reproduce test run with CLI.
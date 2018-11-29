# ocicert: a test framework for OCI distribution certification

A simple test framework for doing OCI distribution certification.

Run the tests like:

```
OCICERT_REGISTRY="docker.io/busybox:latest" make test
```

To run integration tests with a local registry:

```
OCICERT_LOCALREG=1 OCICERT_REGISTRY="127.0.0.1:5000/busybox:latest" make test
```


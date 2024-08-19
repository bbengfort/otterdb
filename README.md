# OtterDB

**A sqlite database replicated using strong consensus for fault tolerance and fast, local reads.**

## Test Cluster

The test cluster is a three replica cluster defined by the docker compose configuration. The replica names are jade, kira, and opal and each replica has its own configuration and data storage in the `tmp` directory.
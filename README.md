# Benchmarks for quic-go

This repository contains benchmarks for [quic-go](https://github.com/lucas-clemente/quic-go/).

## netem setup

netem is only available on Linux. This benchmark suite uses the `tc` command to configure different network conditions, and then runs file transfers using TCP and QUIC.

`tc` has to be run with root priviledges. The following solution is far from optimal, but works fine if you don't care about security (e.g. if you're running the benchmarks in a virtual machine).

1. Add the following lines to */etc/sudoers*:
```
Defaults env_keep += "GOPATH"
Defaults env_keep += "GOROOT"
```
2. Add the file */etc/sudoers.d/%USERNAME%* with the following contents:
```
%USERNAME% ALL=(ALL) NOPASSWD:ALL
```

Tested on Ubuntu 17.10 and 19.04.

# nexus

⛩️ A plugin that can seamlessly link AnyShake products to SeisComP3.

## Getting Started

**Note: This repository is a work in progress (WIP).**

1. Compile the C dynamic library:
    ```bash
    $ make lib
    ```
2. Build the Go application:
    ```bash
    $ CGO_ENABLED=1 go build
    ```
3. Run the `nexus` application by specifying the address of the TCP forwarder service. This address can be obtained from the AnyShake Observer.
    ```bash
    $ export LD_LIBRARY_PATH=$(pwd):$LD_LIBRARY_PATH
    $ ./nexus -address 10.0.0.155:30000
    ```

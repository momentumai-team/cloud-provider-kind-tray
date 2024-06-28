# cloud-provider-kind-tray

Tray application to allow seeing and using load balancers created on kind clusters

## Requirements

- Leverages [cloud-provider-kind](https://github.com/kubernetes-sigs/cloud-provider-kind) to enable kind clusters to have load balancers. A fork of this is installed awaiting acceptance of a [PR](https://github.com/kubernetes-sigs/cloud-provider-kind/pull/99).

## Installation

This will ask for your password to install the application as must create a sudoers file entry to ensure application can be run as root.

```bash
make install
```

## Running

This will run a tray application that you can start and stop to watch for load balancers to be created and ability click on the load balancer to open a browser window on the port

```bash
make run
```

## Credits

[Load Balancer Icon](https://www.flaticon.com/free-icons/load-balancer) provider by [www.flaticon.com](https://www.flaticon.com/)

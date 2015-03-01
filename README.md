# consul-check [![Build Status](https://travis-ci.org/srspnda/consul-check.png)](https://travis-ci.org/srspnda/consul-check)

* Consul website: [consul.io](https://consul.io)
* Consul IRC: `#consul` on Freenode
* Consul mailing list: [Google Groups](https://groups.google.com/group/consul-tool/)

`consul-check` is a CLI helper for executing health checks on systems built
with Consul. Consul is a distributed, highly available, and scalable tool for
service discovery and configuration.

## Usage

Documentation on check definitions is available at the Consul website:

* [Checks](https://consul.io/docs/agent/checks.html)

In general, a check script is free to determine the status of a check in many
ways. This CLI is opinionated, but not declaring itself the only or best way for
implementation. The only limitations of a check definition (placed by Consul
iteself) are the exit codes of a periodically executed script. Specifically:

* Exit code 0 - Check is passing
* Exit code 1 - Check is warning
* Any other code - Check is failing

This is the only convention that Consul depends on. The output of each script
execution is captured and stored in the `notes` field, so that it can be viewed
by a systems administrator.

## Author

* [Justin Poole](mailto:sdpnda@gmail.com)

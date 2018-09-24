# nettools

This is a collection of network-related tools written in Go. It currently includes a version of `netcat`, for linking a TCP socket to stdio and a tool called `netcoupler`, for linking two sockets to each other.

## Usage

```bash
netcat [OPTIONS] addr
netcoupler [OPTIONS] addr1 addr2
```

Both tools have the same options.

- `-l` Listen on the provided address(es) instead of connecting
- `-k` With `-l`, keep the socket(s) open after the remote disconnects

Addresses are in the form `[host]:port`, where `host` can be an IP address or hostname and `port` is the TCP port to listen on or connect to. `host` is only optional in listen mode.
# netcat

This is my implementation of the Unix tool `netcat`. It is much simpler than most other versions, and is not designed to be a drop-in replacement for any of them. It is written in around 100 lines of Go, making it both cross-platform and very easy to maintain.

## Usage

To connect to a network socket: `netcat addr:port`

To listen on a port: `netcat -l :port`

To listen on a port, keeping the socket open even after connections are closed: `netcat -l -k :port`
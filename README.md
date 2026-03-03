# Gopher Status

A simple application to display services on a local computer and allow you to stop and start those services. Originally intended to be run on a Linux computer and provides managers for Docker containers and Systemd services.

For docker containers, the user the `gopherstatus` application is running as will need to be in the `docker` group. Beware that this has security implications you should consider.

For Systemd services, you probably need to add entries to the `sudoers` file to allow it to work.

## Current Status

This is basically feature complete for my needs. At some point, I may revisit this code and add authentication but I did not currently need that.

## How to Deploy

1. Compile the application `go build -o gopherstatus cmd/main.go`.
2. Copy the `gopherstatus` application binary to the location where you will run it from.
3. Copy the `template` and `static` folders to the working directory of the application (probably the same directory you put the `gopherstatus` application binary in).
4. Create a `config.toml` file in the same directory as the `gopherstatus` application binary. See below for an example.

Sample `config.toml`:
```toml
ip = "0.0.0.0"
port = 3000
template_path = "templates"

[status.restreamer]
type = "docker"
container_name = "postgres:18"
image = "postgres"
arguments = [
    "-e",
    "POSTGRES_PASSWORD=yourpasswordhere",
    "-v",
    "/var/postgres:/var/lib/postgresql/data"
]

[status.tailscale]
type = "systemctl"
service_name = "tailscaled"

[status.nginx]
type = "systemctl"
service_name = "nginx"
```
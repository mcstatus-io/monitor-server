# Monitor Server
![](https://img.shields.io/github/languages/code-size/mcstatus-io/monitor-server)
![](https://img.shields.io/github/issues/mcstatus-io/monitor-server)
![](https://img.shields.io/github/actions/workflow/status/mcstatus-io/monitor-server/go.yml)

This is the source code for the monitoring service used internally by mcstatus.io to track the status and uptime information for the dashboard and status pages. If you are looking for the software that retrieves the server statuses itself, please check the [ping-server](https://github.com/mcstatus-io/ping-server) repository instead.

## Requirements

- [Go](https://go.dev/)
- [MongoDB](https://www.mongodb.com/)
- [GNU Make](https://www.gnu.org/software/make/)

## Getting Started

```bash
# 1. Clone the repository (or download from this page)
$ git clone https://github.com/mcstatus-io/monitor-server.git

# 2. Move the working directory into the cloned repository
$ cd monitor-server

# 3. Run the build script
$ make

# 4. Copy the `config.example.yml` file to `config.yml` and modify details as needed
$ cp config.example.yml config.yml

# 5. Start the service
$ ./bin/main
```

## Copyright
&copy; 2022 Jacob Gunther
# slack-ips

[![License: AGPL v3](https://img.shields.io/badge/License-AGPLv3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
![semver](https://img.shields.io/badge/semver-1.0.0-blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/dusnm/slack-ips)](https://goreportcard.com/report/github.com/dusnm/slack-ips)

An HTTP backend for a Slack app that saves bank account numbers to allow for quick sharing with [IPS](https://ips.nbs.rs/en) QR codes,
an instant payment technology in Serbia.

## Screenshots
![Usage in slack](https://github.com/dusnm/slack-ips/blob/main/assets/screenshot-1.png?raw=true)

## Prerequisites
To use this backend, you must [create a Slack app](https://docs.slack.dev/app-management/quickstart-app-settings/) and install it in your workspace.

## Build instructions
You can either build from source, or use prebuilt versions under Docker. A `docker-compose.yml` file is available in 
the docker directory for your convenience.

### Building from source
#### Build requirements:
* A UNIX-like build environment (Windows support is untested, but may work)
* `go >= 1.26`

#### Building
Clone the repository:
```shell
git clone https://github.com/dusnm/slack-ips && cd slack-ips
```

Create the build directory in the source tree:
```shell
mkdir -p build
```
Build either with [`go-task`](https://taskfile.dev/) or manually invoke `go build`:
```shell
task build
```
or
```shell
CGO_ENABLED=0 go build -ldflags='-s -w -extldflags "-static"' -o ./build/slack-ips ./main.go
```

### Docker
You can use prebuilt images from Docker Hub or you can build and tag them yourself using the included `Dockerfile`.
```shell
docker buildx build --platform linux/amd64,linux/arm64 -f ./docker/prod/Dockerfile -t your_organization/slack-ips:latest .
```
## Configuration
Initialize the SQLite database with the included shell script:
```shell
./initdb.sh
```

Use the included `config.example.toml` as a template for configuration.

The application looks for the configuration file in one of these locations and reads values from the first one found.
* `./config.toml` (in the directory where you've placed the compiled binary)
* `$HOME/Library/Application Support/slack-ips/config.toml` (MacOS only)
* `$XDG_CONFIG_HOME/slack-ips/config.toml`

Fill the `[app]` section of the configuration with details such as `bind` and `port` and make special note of the `signing_secret`!

The signing secret, used for URL signing, must be generated using cryptographically secure pseudo-random bytes.
An easy way to do this is to use `/dev/urandom` as a secure source of randomness.
```shell
head -c 20 /dev/urandom | xxd -p
```

Fill the `[slack]` section with details you obtained after creating your Slack app.
The `signing_secret` in this section should be different from the `signing_secret` in the `[app]` section.

Fill the `[db]` section with a `path` to your SQLite database.

## Running the application
If you've configured everything correctly, simply run the application with:
```shell
/path/to/binary/slack-ips
```

## Licensing
This application is free software, licensed under the terms of the GNU Affero General Public License, version 3.

## Remarks
This project contains no AI generated code.

![Developed by a human, not by AI](https://github.com/dusnm/slack-ips/blob/main/assets/no-ai-badge.png?raw=true)
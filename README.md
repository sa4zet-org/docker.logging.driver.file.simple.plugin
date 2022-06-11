[![Check and build](https://github.com/sa4zet-org/docker.logging.driver.file.simple.plugin/actions/workflows/release.yml/badge.svg)](https://github.com/sa4zet-org/docker.logging.driver.file.simple.plugin/actions/workflows/release.yml)
# docker.logging.driver.file.simple.plugin

Logging driver for Docker that saves container raw outputs (stdout and stderr) to a configured directory on the host. If the directory does not exist, the plugin creates it
automatically.

## Install

```bash
docker plugin install --grant-all-permissions ghcr.io/sa4zet-org/docker.logging.driver.file.simple.plugin
```

## Uninstall

```bash
docker plugin disable ghcr.io/sa4zet-org/docker.logging.driver.file.simple.plugin
```

```bash
docker plugin rm ghcr.io/sa4zet-org/docker.logging.driver.file.simple.plugin
```

### Options

All available options are documented here and can be set via `--log-opt="KEY=VALUE"`.

| Key            | Default | Description                 |
|----------------|---------|-----------------------------|
| `log-file-dir` | /tmp/   | File path of the log files. |

## Usage

Run a container using this plugin:

```bash
docker run \
--rm \
--detach \
--name="example_container" \
--log-driver="ghcr.io/sa4zet-org/docker.logging.driver.file.simple.plugin" \
--log-opt="log-file-dir=${HOME}" \
debian:sid-slim \
/bin/sh -c 'ls -lah && >&2 echo "error\nexample\ntext"'
```

# License

https://github.com/sa4zet-org/docker.logging.driver.file.simple.plugin/blob/master/LICENSE

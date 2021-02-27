# docker.logging.driver.file.simple

Logging driver for Docker that saves container raw outputs (stdout and stderr) to a configured directory on the host.

## Install

```
└› docker plugin install sa4zet/docker.logging.driver.file.simple.plugin
Plugin "sa4zet/docker.logging.driver.file.simple.plugin" is requesting the following privileges:
 - network: [host]
 - mount: [/]
Do you grant the above permissions? [y/N] y
latest: Pulling from sa4zet/docker.logging.driver.file.simple.plugin
Digest: sha256:2da1ab07c9f904810ca274da2fca5c96768b3963ab00d7623ff51ccc1a3e4a54
80a804ab35e7: Complete
Installed plugin sa4zet/docker.logging.driver.file.simple.plugin
```

## Check

```
└› docker plugin list
ID             NAME                                                     DESCRIPTION           ENABLED
509f49b713c8   sa4zet/docker.logging.driver.file.simple.plugin   File logging driver   true
```

## Uninstall

```
└› docker plugin disable sa4zet/docker.logging.driver.file.simple.plugin
sa4zet/docker.logging.driver.file.simple.plugin

└› docker plugin rm sa4zet/docker.logging.driver.file.simple.plugin
sa4zet/docker.logging.driver.file.simple.plugin

└› docker plugin list
ID        NAME      DESCRIPTION   ENABLED
```

### Options

All available options are documented here and can be set via `--log-opt="KEY=VALUE"`.

|Key|Default|Description|
|---|---|---|
|`log-file-dir`|/tmp/|File path of the log files.|

## Usage

Run a container using this plugin:

```
└› docker run \
--rm \
--detach \
--name="example_container" \
--log-driver="sa4zet/docker.logging.driver.file.simple.plugin" \
--log-opt="log-file-dir=${HOME}" \
debian:sid-slim \
/bin/sh -c 'ls -lah && >&2 echo "error\nexample\ntext"'
eea0f0e2b87af188465a2c796afce6c5f0d0bbdca56dcd2316cc77284bbdf877

└› cat example_container.out.log
total 72K
drwxr-xr-x   1 root root 4.0K Feb 27 14:58 .
drwxr-xr-x   1 root root 4.0K Feb 27 14:58 ..
-rwxr-xr-x   1 root root    0 Feb 27 14:58 .dockerenv
drwxr-xr-x   2 root root 4.0K Feb  8 00:00 bin
drwxr-xr-x   2 root root 4.0K Jul  9  2019 boot
drwxr-xr-x   5 root root  340 Feb 27 14:59 dev
drwxr-xr-x   1 root root 4.0K Feb 27 14:58 etc
drwxr-xr-x   2 root root 4.0K Jul  9  2019 home
drwxr-xr-x   8 root root 4.0K Feb  8 00:00 lib
drwxr-xr-x   2 root root 4.0K Feb  8 00:00 lib64
drwxr-xr-x   2 root root 4.0K Feb  8 00:00 media
drwxr-xr-x   2 root root 4.0K Feb  8 00:00 mnt
drwxr-xr-x   2 root root 4.0K Feb  8 00:00 opt
dr-xr-xr-x 145 root root    0 Feb 27 14:59 proc
drwx------   2 root root 4.0K Feb  8 00:00 root
drwxr-xr-x   3 root root 4.0K Feb  8 00:00 run
drwxr-xr-x   2 root root 4.0K Feb  8 00:00 sbin
drwxr-xr-x   2 root root 4.0K Feb  8 00:00 srv
dr-xr-xr-x  13 root root    0 Feb 27 14:53 sys
drwxrwxrwt   2 root root 4.0K Feb  8 00:00 tmp
drwxr-xr-x  11 root root 4.0K Feb  8 00:00 usr
drwxr-xr-x  11 root root 4.0K Feb  8 00:00 var

└› cat example_container.err.log
error
example
text
```

# License

https://github.com/sa4zet-org/docker.logging.driver.file.simple/blob/master/LICENSE

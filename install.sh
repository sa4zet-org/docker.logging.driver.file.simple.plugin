#!/bin/bash

docker plugin disable sa4zet/docker.logging.driver.file.simple.plugin
docker plugin rm sa4zet/docker.logging.driver.file.simple.plugin
docker build . --tag sa4zet/docker.logging.driver.file.simple
rm -rf rootfs
mkdir -p rootfs/
id=$(docker create sa4zet/docker.logging.driver.file.simple true)
docker export "${id}" | tar -x -C rootfs
docker rm -f "${id}"
docker plugin create sa4zet/docker.logging.driver.file.simple.plugin .
docker plugin enable sa4zet/docker.logging.driver.file.simple.plugin

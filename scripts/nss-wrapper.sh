#!/bin/sh

# This script should be used as the entrypoint of the Docker container
# to ensure that the passwd file has a reference to the user
# even if the container is started with a non-existant user
# (which is the case by default on OpenShift)

export USER_ID=$(id -u)
export GROUP_ID=$(id -g)
export NSS_WRAPPER_PASSWD=/tmp/passwd
export NSS_WRAPPER_GROUP=/etc/group

cat /etc/passwd > $NSS_WRAPPER_PASSWD
echo "default:x:${USER_ID}:${GROUP_ID}:Default Application User:${HOME}:/sbin/nologin" >> $NSS_WRAPPER_PASSWD

export LD_PRELOAD=libnss_wrapper.so

exec "$@"

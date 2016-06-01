# OpenShift/Git Integration

**Import/Export OpenShift resources from/to a Git repository.**

[![DockerHub](https://img.shields.io/badge/docker-vbehar%2Fopenshift--git-008bb8.svg)](https://hub.docker.com/r/vbehar/openshift-git/)
[![Travis](https://travis-ci.org/vbehar/openshift-git.svg?branch=master)](https://travis-ci.org/vbehar/openshift-git)
[![Circle CI](https://circleci.com/gh/vbehar/openshift-git/tree/master.svg?style=svg)](https://circleci.com/gh/vbehar/openshift-git/tree/master)

`openshift-git` is an [OpenShift](https://www.openshift.org/) client that can be used to export resources from your cluster to a Git repository, or to import resources from a Git repository to your cluster. It is written in [Go](https://golang.org/) so it just produces a single binary.

## Goal

The main goal of this project is to store your OpenShift resources (buildconfigs, deploymentconfigs, ...) in a Git repository, so that you can record every change, and thus have an easy access to an older version, thanks to Git history.

While it can't really be used as an audit tool (it won't store who did the change, or why they did it), it will still record what has been changed and when, which is quite useful.

It can be used to export either a single namespace (so that if you are not a cluster-admin, you can still benefit from it), or the whole cluster (obviously only if you are a cluster-admin).

## Features

* The main feature is the **daemon export** mode, in which `openshift-git` will run forever, and commit to the Git repository every change that happens in the cluster.
* But it can also be used as a one-time export, if you prefer periodic exports.
* An import command is planned, but not yet implemented.

## Usage

Get the binary from the latest release, then just run

```
openshift-git
```

and it will print the available commands, options, and some examples.

## How It Works

The `export` command has 2 modes:

* the standard export, that will list all requested resources, save them to the filesystem and then commit them to the Git repository.
* the daemon export, that will start by listing all requested resources, and then open a "watch" to listen for every change, and commit them to the Git repository.

By default it will only commit to the local Git repository, but if you provide the URL of a remote Git repository, it will periodically push the local commits to the remote repository.

It can export as little or as many different types of resources as you need, depending on how you start it.

## Running on OpenShift

There are 2 ways to deploy this application on an OpenShift cluster:

* For exporting resources from the whole cluster (requires cluster-admin role):

  * create a `git-exporter` service account:
  
  ```
  oc create serviceaccount git-exporter
  ```

  * give the `cluster-reader` role to the newly created service account:
  
  ```
  oc adm policy add-cluster-role-to-user cluster-reader system:serviceaccount:$(oc project -q):git-exporter
  ```

  * if you want to push to a remote git repository, you need to create a secret for your SSH key:
  
  ```
  oc create secret generic mysshkey --from-file=publickey=$HOME/.ssh/id_rsa.pub --from-file=privatekey=$HOME/.ssh/id_rsa --from-file=config=$HOME/.kube/ssh-config
  ```
  
  With the following content in the `$HOME/.kube/ssh-config` file:
  
  ```
  Host *
  IdentityFile ~/.ssh/privatekey
  StrictHostKeyChecking no
  ```

  * create a new application from the provided [openshift-template-full-cluster.yml](openshift-template-full-cluster.yml) template, and overwrite some parameters:

  ```
  oc new-app -f openshift-template-full-cluster.yml -p SERVICE_ACCOUNT=git-exporter,SSH_KEYS_SECRET=mysshkey,REMOTE_GIT_REPOSITORY_URL=git@github.com:USER/REPO.git
  ```

* For exporting resources from a single namespace (does not requires specific rights):

  * if you want to push to a remote git repository, you need to create a secret for your SSH key:
  
  ```
  oc create secret generic mysshkey --from-file=publickey=$HOME/.ssh/id_rsa.pub --from-file=privatekey=$HOME/.ssh/id_rsa --from-file=config=$HOME/.kube/ssh-config
  ```
  
  With the following content in the `$HOME/.kube/ssh-config` file:
  
  ```
  Host *
  IdentityFile ~/.ssh/privatekey
  StrictHostKeyChecking no
  ```

  * create a new application from the provided [openshift-template-single-namespace.yml](openshift-template-single-namespace.yml) template, and overwrite some parameters:

  ```
  oc new-app -f openshift-template-single-namespace.yml -p SSH_KEYS_SECRET=mysshkey,REMOTE_GIT_REPOSITORY_URL=git@github.com:USER/REPO.git
  ```

## Running locally

If you want to run it on your laptop:

* Install [Go](http://golang.org/) (tested with 1.6) and [setup your GOPATH](https://golang.org/doc/code.html)
* clone the sources in your `GOPATH`

  ```
  git clone https://github.com/vbehar/openshift-git.git $GOPATH/src/github.com/vbehar/openshift-git
  ```

* install [godep](https://github.com/tools/godep) (to use the vendored dependencies)

  ```
  go get github.com/tools/godep
  ```

* build the binary with godep:

  ```
  cd $GOPATH/src/github.com/vbehar/openshift-git
  godep go install
  ```

* and run it:

  ```
  $GOPATH/bin/openshift-git
  ```

* enjoy!

## License

Copyright 2016 the original author or authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

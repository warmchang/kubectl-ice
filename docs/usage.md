# Contents
* [Introduction](#introduction)
* [Usage](#usage)
* [Flags](#flags)
* [Command](#command)
  * [Examples](#example)
* [CPU](#cpu)
  * [Examples](#example-1)
* [Environment](#environment)
  * [Examples](#example-2)
* [Image](#image)
  * [Examples](#example-3)
* [IP](#ip)
  * [Examples](#example-4)
* [Memory](#memory)
  * [Examples](#example-5)
* [Ports](#ports)
  * [Examples](#example-6)
* [Probes](#probes)
  * [Examples](#example-7)
* [Restarts](#restarts)
  * [Examples](#example-8)
* [Status](#status)
  * [Examples](#example-9)
* [Volumes](#volumes)
  * [Examples](#example-10)

## Introduction
A kubectl plugin that lets you can see the running configuration of all containers
 that are running inside pods, I created it so I could peer inside the pods and see
 the details of containers (sidecars) running in a pod and then extended it so all
 containers could be viewed at once.   

ice lists useful information about the sidecar containers present inside a
 pod, useful for trouble shooting multi container issues you can view volume, 
 image, port and executable configurations, along with current cpu and memory
  metrics all at the container level (requires metrics server)

## Usage
ice usage is split in to sub commands with each following commands are available for `kubectl ice`
```
kubectl ice command    # retrieves the command line and any arguments specified at the container level
kubectl ice cpu        # return cpu requests size, limits and usage of each container
kubectl ice help       # Shows the help screen
kubectl ice image      # list the image name and pull status for each container
kubectl ice ip         # list ip addresses of all pods in the namespace listed
kubectl ice memory     # return memory requests size, limits and usage of each container
kubectl ice ports      # shows ports exposed by the containers in a pod
kubectl ice probes     # shows details of configured startup, readiness and liveness probes of each container
kubectl ice restarts   # show restart counts for each container in a named pod
kubectl ice status     # list status of each container in a pod
kubectl ice volumes    # list all container volumes with mount points
```

## Flags
All standard kubectl flags are supported including the beow, see the examples section of each command for usage info:
```
  -A, --all-namespaces                 List containers from pods in all namespaces
  -c, --container string               Container name. If set shows only the named containers containers in the pod
      --context string                 The name of the kubeconfig context to use
      --match string                   excludes results, comma seperated list of COLUMN OP VALUE, where OP can be one of ==,<,>,<=,>= and != 
  -n, --namespace string               If present, the namespace scope for this CLI request
  -l, --selector string                Selector (label query) to filter on
```
selected subcommands also support the following flags
```
  -p, --previous         show previous state
  -r, --raw              show raw uncooked values
      --sort string      Sort by column
      --oddities         show only the outlier rows that dont fall within the computed range (requires min 5 rows in output)
```

## Command
Retrieves the command line and any arguments specified at the container level

Prints command and arguments used to start each container (if specified), single pods and
containers can be selected by name.  If no name is specified the container commands of all pods
in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice command [flags]

Aliases:
  command, cmd, exec, args

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers command info from pods
kubectl-ice command

# List container command info from pods output in JSON format
kubectl-ice command -o json

# List container command info from a single pod
kubectl-ice command my-pod-4jh36

# List command info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice command -c web-container

# List command info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice command -c web-container --sort '!CONTAINER'

# List command info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice command -c web-container --sort PODNAME

# List container command info from all pods where label app matches web
kubectl-ice command -l app=web

# List container command info from all pods where the pod label app is either web or mail
kubectl-ice command -l "app in (web,mail)"
```
## CPU
Show configured cpu size, limit and % usage of each container

Prints the current cpu usage along with configured requests and limits. The calculated % fields
serve as an easy way to see how close you are to the configured sizes.  By specifying the -r
flag you can see raw unfiltered values.  If no name is specified the container cpu details
of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice cpu [flags]

Flags:
      --oddities                       show only the outlier rows that dont fall within the computed range
  -r, --raw                            show raw values

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers cpu info from pods
kubectl-ice cpu

# List container cpu info from pods output in JSON format
kubectl-ice cpu -o json

# List container cpu info from a single pod
kubectl-ice cpu my-pod-4jh36

# List cpu info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice cpu -c web-container

# List cpu info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice cpu -c web-container --sort '!CONTAINER'

# List cpu info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice cpu -c web-container --sort PODNAME

# List container cpu info from all pods where label app matches web
kubectl-ice cpu -l app=web

# List container cpu info from all pods where the pod label app is either web or mail
kubectl-ice cpu -l "app in (web,mail)"
```
## Environment
List the env name and value for each container

Print the the environment variables used in running containers in a pod, single pods
and containers can be selected by name. If no name is specified the environment details of all pods in
the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice environment [flags]

Aliases:
  environment, env, vars

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers env info from pods
kubectl-ice env

# List container env info from pods output in JSON format
kubectl-ice env -o json

# List container env info from a single pod
kubectl-ice env my-pod-4jh36

# List env info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice env -c web-container

# List env info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice env -c web-container --sort '!CONTAINER'

# List env info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice env -c web-container --sort PODNAME

# List container env info from all pods where label app matches web
kubectl-ice env -l app=web

# List container env info from all pods where the pod label app is either web or mail
kubectl-ice env -l "app in (web,mail)"
```
## Image
List the image name and pull status for each container

Print the the image used for running containers in a pod including the pull policy, single pods
and containers can be selected by name. If no name is specified the image details of all pods in
the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice image [flags]

Aliases:
  image, im

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers image info from pods
kubectl-ice image

# List container image info from pods output in JSON format
kubectl-ice image -o json

# List container image info from a single pod
kubectl-ice image my-pod-4jh36

# List image info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice image -c web-container

# List image info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice image -c web-container --sort '!CONTAINER'

# List image info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice image -c web-container --sort PODNAME

# List container image info from all pods where label app matches web
kubectl-ice image -l app=web

# List container image info from all pods where the pod label app is either web or mail
kubectl-ice image -l "app in (web,mail)"
```
## IP
List ip addresses of all pods in the namespace listed

Prints the known IP addresses of the specified pod(s). if no pod is specified the IP address of
all pods in the current namespace are shown.

``` shell
Usage:
  kubectl-ice ip [flags]

```
also includes standard common kubectl flags

### Examples
``` shell
# List IP address of pods
kubectl-ice ip

# List IP address of pods output in JSON format
kubectl-ice ip -o json

# List IP address a single pod
kubectl-ice ip my-pod-4jh36

# List IP address of all pods where label app matches web
kubectl-ice ip -l app=web

# List IP address of all pods where the pod label app is either web or mail
kubectl-ice ip -l "app in (web,mail)"
```
## Memory
Show configured memory size, limit and % usage of each container

Prints the current memory usage along with configured requests and limits. The calculated % fields
serve as an easy way to see how close you are to the configured sizes.  By specifying the -r
flag you can see raw unfiltered values.  If no name is specified the container memory details
of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice memory [flags]

Aliases:
  memory, mem

Flags:
      --oddities                       show only the outlier rows that dont fall within the computed range
  -r, --raw                            show raw values

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers memory info from pods
kubectl-ice memory

# List container memory info from pods output in JSON format
kubectl-ice memory -o json

# List container memory info from a single pod
kubectl-ice memory my-pod-4jh36

# List memory info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice memory -c web-container

# List memory info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice memory -c web-container --sort '!CONTAINER'

# List memory info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice memory -c web-container --sort PODNAME

# List container memory info from all pods where label app matches web
kubectl-ice memory -l app=web

# List container memory info from all pods where the pod label app is either web or mail
kubectl-ice memory -l "app in (web,mail)"
```
## Ports
Shows ports exposed by the containers in a pod

Print a details of service ports exposed by containers in a pod. Details include the container
name, port number and protocol type. Port name and host port are only show if avaliable. If no
name is specified the container port details of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice ports [flags]

Aliases:
  ports, port, po

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers port info from pods
kubectl-ice ports

# List container port info from pods output in JSON format
kubectl-ice ports -o json

# List container port info from a single pod
kubectl-ice ports my-pod-4jh36

# List port info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice ports -c web-container

# List port info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice ports -c web-container --sort '!CONTAINER'

# List port info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice ports -c web-container --sort PODNAME

# List container port info from all pods where label app matches web
kubectl-ice ports -l app=web

# List container port info from all pods where the pod label app is either web or mail
kubectl-ice ports -l "app in (web,mail)"
```
## Probes
Shows details of configured startup, readiness and liveness probes of each container

Prints details of the currently configured startup, liveness and rediness probes for each
container. Details like the delay timeout and action are printed along with the configured probe
type. If no name is specified the container probe details of all pods in the current namespace
are shown.

``` shell
Usage:
  kubectl-ice probes [flags]

Aliases:
  probes, probe

```
also includes standard common kubectl flags

### Examples
``` shell
# List containers probe info from pods
kubectl-ice probes

# List container probe info from pods output in JSON format
kubectl-ice probes -o json

# List container probe info from a single pod
kubectl-ice probes my-pod-4jh36

# List probe info for all containers named web-container searching all
# pods in the current namespace
kubectl-ice probes -c web-container

# List probe info for all containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice probes -c web-container --sort '!CONTAINER'

# List probe info for all containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice probes -c web-container --sort PODNAME

# List container probe info from all pods where label app matches web
kubectl-ice probes -l app=web

# List container probe info from all pods where the pod label app is either web or mail
kubectl-ice probes -l "app in (web,mail)"
```
## Restarts
Show restart counts for each container in a named pod

Prints container name and restart count for individual containers. If no name is specified the
container restart counts of all pods in the current namespace are shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice restarts [flags]

Aliases:
  restarts, restart

Flags:
      --oddities                       show only the outlier rows that dont fall within the computed range

```
also includes standard common kubectl flags

### Examples
``` shell
# List individual container restart count from pods
kubectl-ice restarts

# List conttainers restart count from pods output in JSON format
kubectl-ice restarts -o json

# List restart count from all containers in a single pod
kubectl-ice restarts my-pod-4jh36

# List restart count of all containers named web-container searching all
# pods in the current namespace
kubectl-ice restarts -c web-container

# List restart count of containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice restarts -c web-container --sort '!CONTAINER'

# List restart count of containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice restarts -c web-container --sort PODNAME

# List container restart count from all pods where label app equals web
kubectl-ice restarts -l app=web

# List restart count from all containers where the pod label app is either web or mail
kubectl-ice restarts -l "app in (web,mail)"
```
## Status
List status of each container in a pod

Prints container status information from pods, current and previous exit code, reason and signal
are shown slong with current ready and running state. Pods and containers can also be selected
by name. If no name is specified the container state of all pods in the current namespace is
shown.

The T column in the table output denotes S for Standard and I for init containers

``` shell
Usage:
  kubectl-ice status [flags]

Aliases:
  status, st

Flags:
      --oddities                       show only the outlier rows that dont fall within the computed range
  -p, --previous                       show previous state

```
also includes standard common kubectl flags

### Examples
``` shell
# List individual container status from pods
kubectl-ice status

# List conttainers status from pods output in JSON format
kubectl-ice status -o json

# List status from all container in a single pod
kubectl-ice status my-pod-4jh36

# List previous container status from a single pod
kubectl-ice status -p my-pod-4jh36

# List status of all containers named web-container searching all
# pods in the current namespace
kubectl-ice status -c web-container

# List status of containers called web-container searching all pods in current
# namespace sorted by container name in descending order (notice the ! charator)
kubectl-ice status -c web-container --sort '!CONTAINER'

# List status of containers called web-container searching all pods in current
# namespace sorted by pod name in ascending order
kubectl-ice status -c web-container --sort PODNAME

# List container status from all pods where label app equals web
kubectl-ice status -l app=web

# List status from all containers where the pods label app is either web or mail
kubectl-ice status -l "app in (web,mail)"
```
## Volumes
Display container volumes and mount points

Prints configured volume information at the container level, volume type, backing information,
read-write state and mount point are all avaliable, volume size is only available if found in
the pod configuration. If no name is specified the volume information for all pods in the
current namespace are shown.

``` shell
Usage:
  kubectl-ice volumes [flags]

Aliases:
  volumes, volume, vol

```
also includes standard common kubectl flags

### Examples
``` shell
# List volumes from containers inside pods from current namespace
kubectl-ice volumes

# List volumes from conttainers output in JSON format
kubectl-ice volumes -o json

# List all container volumes from a single pod
kubectl-ice volumes my-pod-4jh36

# List volumes from all containers named web-container searching all
# pods in the current namespace
kubectl-ice volumes -c web-container

# List volumes from container web-container searching all pods in current
# namespace sorted by volume name in descending order (notice the ! charator)
kubectl-ice volumes -c web-container --sort '!VOLUME'

# List volumes from container web-container searching all pods in current
# namespace sorted by volume name in ascending order
kubectl-ice volumes -c web-container --sort MOUNT-POINT

# List container volume info from all pods where label app equals web
kubectl-ice volumes -l app=web

# List volumes from all containers where the pod label app is web or mail
kubectl-ice volumes -l "app in (web,mail)"
```

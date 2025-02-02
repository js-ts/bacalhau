# Running locally with the 'devstack' command

The `devstack` command of bacalhau will start a 3 node cluster alongside isolated ipfs servers.

This is useful to kick the tires and/or developing on the codebase.  It's also the tool used by some of the tests.

## Pre-requisites

 * x86_64 linux host
    * Ubuntu 20.0+ has most often been used for development and testing
    * Note: Mac M1 (ARM64) compatible builds are not yet supported at this time. Please consider development on a hosted alternative, such as [Gitpod](https://gitpod.io/#https://github.com/filecoin-project/bacalhau)
 * Go >= 1.17
 * IPFS v0.11
 * [Docker Engine](https://docs.docker.com/get-docker/)
 * A build of the [latest Bacalhau release](https://github.com/filecoin-project/bacalhau/releases/)

## (Optional) Building Bacalhau from source

```
sudo apt-get update && sudo apt-get install -y make gcc zip
sudo snap install go --classic
wget https://github.com/filecoin-project/bacalhau/archive/refs/heads/main.zip
unzip main.zip
cd bacalhau-main
go build

```



## Start the cluster

```bash
make devstack
```

This will start a 3 node bacalhau cluster connected with libp2p.

Each node has it's own ipfs server isolated using the `IPFS_PATH` environment variable and it's own API RPC server isolated using a random port.

Once everything has started up - you will see output like the following:

```bash
-------------------------------
node 0
-------------------------------

export IPFS_API_PORT_0=36825
export IPFS_PATH_0=/tmp/bacalhau-ipfs-devstack4061398535
export API_PORT_0=43079

cid=$(IPFS_PATH=/tmp/bacalhau-ipfs-devstack4061398535 ipfs add -q testdata/grep_file.txt)
curl -XPOST http://127.0.0.1:36825/api/v0/id

-------------------------------
node 1
-------------------------------

export IPFS_API_PORT_1=46023
export IPFS_PATH_1=/tmp/bacalhau-ipfs-devstack3414455455
export API_PORT_1=43079

cid=$(IPFS_PATH=/tmp/bacalhau-ipfs-devstack3414455455 ipfs add -q testdata/grep_file.txt)
curl -XPOST http://127.0.0.1:46023/api/v0/id

-------------------------------
node 2
-------------------------------

export IPFS_API_PORT_2=40277
export IPFS_PATH_2=/tmp/bacalhau-ipfs-devstack2766996210
export API_PORT_2=43079

cid=$(IPFS_PATH=/tmp/bacalhau-ipfs-devstack2766996210 ipfs add -q testdata/grep_file.txt)
curl -XPOST http://127.0.0.1:40277/api/v0/id
```

## New Terminal Window
* Open an additional terminal window to be used for data submission to the local IPFS instances and and job submission to the 3 node devestack Bacalhau cluster.
* Copy and paste the IPFS and API port variables into the new terminal window.

## Add files to IPFS

Each node has it's own `IPFS_PATH` value which points to a path on the local filesystem.  This allows to use the ipfs cli to test adding files to one or multiple nodes.  This is especially useful when you want to test self selection of a job based on whether the cid is *local* to that node.

To add a file to only one of ipfs node within the devstack cluster, execute the `ipfs add` in the following manner:

```bash
cid=$( IPFS_PATH=$IPFS_PATH_0 ipfs add -q ./testdata/grep_file.txt )
```
*Note: the CID is saved as an environment variable so that it can be referenced in the job submission step.

## Set a json rpc port

Each node has it's own `--api-port` value.  This means you can use the `go run .` cli in isolation from the other 2 nodes.

For example - to view the current job list from the perspective of only one of the 3 nodes:

```bash
# Note: replace 12345 this with the correct port from the output
go run . --api-port=$API_PORT_0 --api-host=localhost list
```

## Submit a simple job

This will submit a simple job to a single node:

```bash
cid=$( IPFS_PATH=$IPFS_PATH_0 ipfs add -q ./testdata/grep_file.txt )
go run . --api-port=$API_PORT_0 --api-host=localhost run -v $cid:/file.txt ubuntu grep kiwi /file.txt
go run . --api-port=$API_PORT_0 --api-host=localhost list --wide
```

After a short while - the job should be in `complete` state.

```
kai@xwing:~/projects/bacalhau$ go run . --api-port=$API_PORT_0 --api-host=localhost list --wide
 ID        JOB                                INPUTS  OUTPUTS  CONCURRENCY  NODE      STATE     RESULT                                               
 22b53c20  docker ubuntu grep kiwi /file.txt       1        0            1  QmedX1zE  complete  /ipfs/QmYLFuXZv8h1Bc1cArbs5VXrE4o5hE4tVh55iqtjQWoDtW 
```

We can see the results have been written back to ipfs for us - let's copy the result path to a variable

Copy the job id into a variable:

```bash
RESULT_PATH=/ipfs/QmYLFuXZv8h1Bc1cArbs5VXrE4o5hE4tVh55iqtjQWoDtW 
```

Now we can view the results:

```bash
IPFS_PATH=$IPFS_PATH_0 ipfs ls $RESULT_PATH
IPFS_PATH=$IPFS_PATH_0 ipfs cat $RESULT_PATH/stdout
```

## run 3 node job

Now let's run a job across all 3 nodes.  To do this, we need to add the cid to all the IPFS servers so the job will be selected to run across all 3 nodes:

```bash
cid=$( IPFS_PATH=$IPFS_PATH_0 ipfs add -q ./testdata/grep_file.txt )
IPFS_PATH=$IPFS_PATH_1 ipfs add -q ./testdata/grep_file.txt
IPFS_PATH=$IPFS_PATH_2 ipfs add -q ./testdata/grep_file.txt
```

Then we submit the job but with `--concurrency` setting:

```bash
go run . --api-port=$API_PORT_0 --api-host=localhost run --concurrency=3 -v $cid:/file.txt ubuntu grep kiwi /file.txt
go run . --api-port=$API_PORT_0 --api-host=localhost list --wide
```

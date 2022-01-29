#!/bin/sh -l

export cluster_url=$1
export token=$2
export configmap_name=$3
export namespace=$4
export protocols=$5
export config_name=$6

/endpoint-discover
#!/bin/bash
clear

dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd $dir
cd ../../../..
root=$(pwd)


cd $root/services/security/docker

sudo docker-compose up $1


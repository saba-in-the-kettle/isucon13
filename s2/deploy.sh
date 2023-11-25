#! /bin/bash

set -eu

cd "$(dirname "$0")"

server_name=$(basename `pwd`)

echo "Deploying $server_name"
rsync -azv --no-owner --no-group --inplace ./fs/ root@$server_name:/

cat ./post-deploy.sh | ssh $server_name bash -

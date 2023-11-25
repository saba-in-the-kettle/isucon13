#! /bin/bash

set -eu

server_name=$1

mkdir -p $server_name/fs/etc

rsync -azv --no-owner --no-group root@$server_name:/etc/mysql ./$server_name/fs/etc/
rsync -azv --no-owner --no-group root@$server_name:/etc/nginx ./$server_name/fs/etc/
find  ./$server_name/fs -type l | xargs -I {} unlink {}

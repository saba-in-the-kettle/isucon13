set -x

sudo chmod 777 /var/log/nginx
sudo chmod 777 /var/log/nginx/*
sudo chmod 777 /var/log/mysql
sudo chmod 777 /var/log/mysql/*

set -e

truncate -s 0 /var/log/mysql/mysql-slow.log
truncate -s 0 /var/log/nginx/access.log

sudo systemctl disable --now nginx
sudo systemctl disable --now mysql

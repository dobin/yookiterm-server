#!/bin/bash

cd /home/yookiterm/yookiterm-challenges
git pull

cd /home/yookiterm/yookiterm-slides
git pull

cd /home/yookiterm/yookiterm
git pull

cd /home/yookiterm/yookiterm-server
git pull
/usr/local/go/bin/go build

sudo /bin/systemctl restart yookiterm

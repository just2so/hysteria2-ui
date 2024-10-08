#!/bin/bash

clear

systemctl stop h-ui

rm -rf /etc/systemd/system/h-ui.service /usr/local/h-ui/

echo "卸载成功!!!"

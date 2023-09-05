#!/bin/bash

cp OBIII.service /etc/systemd/system/OBIII.service
systemctl daemon-reload
systemctl stop OBIII
systemctl enable --now OBIII

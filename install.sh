#!/bin/bash
go env > OBIII.env
cp OBIII.service /etc/systemd/system/OBIII.service
systemctl daemon-reload
systemctl stop OBIII
systemctl enable --now OBIII

#!/bin/bash
killall gmgo-admin # kill go-admin service
echo "stop gmgo-admin success"
ps -aux | grep gmgo-admin
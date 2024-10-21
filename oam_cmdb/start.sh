#!/bin/bash

BASE_DIR=`cd $(dirname $0); pwd`
echo ${BASE_DIR}
nohup ${BASE_DIR}/oam >>${BASE_DIR}/nohup.out 2>&1 &
pause

#!/bin/bash
if [ -z $CI_HOME ];
then
   CI_HOME="/root/cid"
fi

if [ -z $CI_TOKENNAME ];
then
   CI_TOKENNAME="stake"
fi

CI_LOGDIR=$CI_HOME/logs
if [ ! -d $CI_LOGDIR ]; then
    mkdir -p $CI_LOGDIR
fi

/opt/cli-linux rest-server --laddr=tcp://0.0.0.0:80 >> $CI_LOGDIR/rest-output.log 2>&1
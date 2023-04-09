#!/bin/bash
CURR=$(pwd)
DOCF=$CURR/build/Dockerfile

echo "##teamcity[blockOpened name='BuildBinary']"
cd $CURR/cmd/hideme/ && go build -o ../../build/bin/hideme || exit 1
echo "##teamcity[blockClosed name='BuildBinary']"

echo "##teamcity[blockOpened name='BuildDockerfile']"
echo "FROM ubuntu:20.04" | tee $DOCF
echo "ADD ./bin /usr/local/hideme" | tee -a $DOCF
echo "WORKDIR /usr/local/hideme" | tee -a $DOCF
echo "CMD ./hideme server --port=8095" | tee -a $DOCF
echo "##teamcity[blockClosed name='BuildDockerfile']"

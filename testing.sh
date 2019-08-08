#!/bin/bash

echo "This script will step through a number of tests agains plunder to ensure that functionality is as expected"
echo "Building plunder with [go build]"

INSECURE="-k"
PLUNDERURL="https://localhost:60443"

go build

retVal=$?
if [ $retVal -ne 0 ]; then
    echo "Error at go build"
    exit
fi

echo "Check for no version information"
v=$(./plunder version | grep Version | awk '{ print $2 }')
rm plunder

if [ -z "$v" ]
then
      echo "Version is empty"
else
      echo "Version is NOT empty"
fi

echo "Building plunder with [make build]"

make build

retVal=$?
if [ $retVal -ne 0 ]; then
    echo "Error at make build"
    exit
fi

echo "Check for version information"
v=$(./plunder version | grep Version | awk '{ print $2 }')
if [ -z "$v" ]
then
      echo "Version is empty"
else
      echo "Version is NOT empty [$v]"
fi

echo "Plunder server configuration, temporary output will live in ./testing"
mkdir testing
./plunder config server -p > ./testing/server_test_config.json
./plunder config deployment -p > ./testing/deployment_config.json
./plunder config server -o yaml > ./testing/server_test_config.yaml
./plunder config deployment -o yaml > ./testing/deployment_config.yaml
echo "Generating API Server certificates in ~./plunderserver.yaml"
./plunder config apiserver server

echo "Creating alternative configuration with services enabled"
sed '/enableHTTP/s/false/true/' ./testing/server_test_config.json > ./testing/server_test_http_config.json


echo "Examining detected configuration"
echo "Checking for Adapter"
v=$(grep adapter ./testing/server_test_config.json | awk ' {print $2 }' | tr -d '"' | tr -d ',')
if [ -z "$v" ]
then
      echo "Adapter is empty"
else
      echo "Adapter is NOT empty [$v]"
fi

echo "Checking for Gateway Address"
v=$(grep gatewayDHCP ./testing/server_test_config.json | awk ' {print $2 }' | tr -d '"' | tr -d ',')
if [ -z "$v" ]
then
      echo "Gateway is empty"
else
      echo "Gateway is NOT empty [$v]"
fi

i=$(id -u)
n=$(id -un)
if [[ $i -gt 0 ]]
then
      echo "Testing as current user [NAME = $n / ID = $i]"
      echo "Starting with disabled configuration"
      ./plunder server --config ./testing/server_test_config.json
      retVal=$?
      if [ $retVal -ne 0 ]; then
          echo "Plunder correctly didn't start"
      fi
      echo "Starting with enabled HTTP configuration (check OSX)"
      sudo ./plunder server --config ./testing/server_test_http_config.json &
      retVal=$?
      if [ $retVal -ne 0 ]; then
          echo "Plunder correctly didn't start"
          exit 1
      fi
      echo "Sleeping for 3 seconds to ensure plunder has started"
      sleep 3
      echo "Print Configuration info"; echo "--------------------------"
      curl $INSECURE $PLUNDERURL/config; echo ""
      echo "Print Deployments info"; echo "--------------------------"
      curl $INSECURE $PLUNDERURL/deployments; echo ""
      echo "POST JSON Deployment to Plunder API"
      curl $INSECURE -X POST -d "@./testing/deployment_config.json" $PLUNDERURL/deployments
      echo "Print (UPDATED) Deployment info"; echo "--------------------------"
      curl $INSECURE $PLUNDERURL/deployments; echo ""
      echo "POST YAML Deployment to Plunder API"
      curl $INSECURE -X POST --data-binary "@./testing/deployment_config.yaml" $PLUNDERURL/deployments -H "Content-type: text/x-yaml"
      echo "Print (UPDATED) Deployment info"; echo "--------------------------"
      curl $INSECURE $PLUNDERURL/deployments; echo ""
      sudo kill -9 $( ps -ef | grep -i plunder | grep -v -e 'sudo' -e 'grep' | awk '{ print $2 }')     
      wait $! 2>/dev/null
      sleep 1
else 
      echo "Skipping permission tests as running as root"
fi

echo "The following tests rely on sudo, with NOPASSWD enabled"

echo "Starting with disabled configuration"

retVal=$?
if [ $retVal -ne 0 ]; then
    echo "Error at make build"
    exit
fi

echo "To remote [./testing/] directory, and [./plunder] binary"
echo "rm -rf ./testing/ ./plunder"
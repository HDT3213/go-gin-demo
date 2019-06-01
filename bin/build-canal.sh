#!/usr/bin/env bash

ROOT_DIR=$(cd "$(dirname "$0")"; cd ..; pwd)
cd ${ROOT_DIR}/app/main/canal
go build -o ${ROOT_DIR}/target/canal
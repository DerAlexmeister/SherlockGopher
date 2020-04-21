#!/bin/bash
# Bashscript to install seabolt for ubuntu 18.04

apt install pkg-config
wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb
apt install -f ./seabolt-1.7.4-Linux-ubuntu-18.04.deb
rm ./seabolt-1.7.4-Linux-ubuntu-18.04.deb
export PKG_CONFIG_PATH="/usr/local/share/pkgconfig"
export PKG_CONFIG_PATH=/seabolt/build/dist/share/pkgconfig
export LD_LIBRARY_PATH=/seabolt/build/dist/lib64
export C_INCLUDE_PATH=/seabolt/build/dist/includ
go get github.com/neo4j/neo4j-go-driver/neo4j
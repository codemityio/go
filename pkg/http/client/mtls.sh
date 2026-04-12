#!/bin/bash

# Guide: https://www.golinuxcloud.com/mutual-tls-authentication-mtls/

set -e

DIR=testdata/mtls

mkdir -p "${DIR}"

reset() {
  rm -rf \
    "${DIR}/01.pem" \
    "${DIR}/02.pem" \
    "${DIR}/serial" \
    "${DIR}/serial.old" \
    "${DIR}/index.txt" \
    "${DIR}/index.txt.old" \
    "${DIR}/index.txt.attr" \
    "${DIR}/index.txt.attr.old"
}

prepare() {
  touch "${DIR}/index.txt"
  echo 01 > "${DIR}/serial"
}

# Certificate Authority Certificate

openssl genrsa -out "${DIR}/ca.key" 4096 2>/dev/null

openssl req -new -x509 \
  -config testdata/mtls.conf \
  -days 3650 \
  -subj "/C=GB/ST=England/L=London/O=TestOrganisation/OU=Information Technology/CN=TestOrganisation" \
  -key "${DIR}/ca.key" \
  -out "${DIR}/ca.crt" 2>/dev/null

# Client Certificate

reset
prepare

openssl genrsa -out "${DIR}/client.key" 4096 2>/dev/null

openssl req -new \
  -subj "/C=GB/ST=England/L=London/O=TestOrganisation/OU=Information Technology/CN=TestOrganisation" \
  -key "${DIR}/client.key" \
  -out "${DIR}/client.csr" 2>/dev/null

openssl ca -notext -batch \
  -config testdata/mtls.conf \
  -extfile testdata/client-ext.conf \
  -days 1650 \
  -in "${DIR}/client.csr" \
  -out "${DIR}/client.crt" 2>/dev/null

# Server Certificate

reset
prepare

openssl genrsa -out "${DIR}/server.key" 4096 2>/dev/null

openssl req -new \
  -subj "/C=GB/ST=England/L=London/O=TestOrganisation/OU=Information Technology/CN=TestOrganisation" \
  -key "${DIR}/server.key" \
  -out "${DIR}/server.csr" 2>/dev/null

openssl ca -notext -batch \
  -config testdata/mtls.conf \
  -extfile testdata/server-ext.conf \
  -days 1650 \
  -in "${DIR}/server.csr" \
  -out "${DIR}/server.crt" 2>/dev/null

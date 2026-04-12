#!/bin/bash

mkdir -p testdata
openssl genpkey -algorithm RSA -out testdata/private.pem 2>/dev/null
openssl rsa -in testdata/private.pem -pubout -out testdata/public.pem 2>/dev/null
openssl rsa -in testdata/private.pem -out testdata/rsa.pem 2>/dev/null

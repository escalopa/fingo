#!/bin/bash

function init() {
  sleep 3
  cockroach init --insecure 2>>/dev/null
  cockroach sql --insecure --execute "CREATE DATABASE IF NOT EXISTS fingo_wallet;" >>/dev/null
}

init &
cockroach start --insecure --join=crdb --listen-addr=0.0.0.0:26257 --advertise-addr=crdb

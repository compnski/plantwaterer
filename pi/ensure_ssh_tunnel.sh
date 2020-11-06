#!/usr/bin/env bash
#set -euo pipefail

createTunnel() {
  /usr/bin/ssh -N -R 2222:localhost:22 pi-bastion
  if [[ $? -eq 0 ]]; then
    echo Tunnel to jumpbox created successfully on `date`
  else
    echo An error occurred creating a tunnel to jumpbox. RC was $? on `date`
  fi
}
/bin/pidof ssh
if [[ $? -ne 0 ]]; then
  echo Creating new tunnel connection on `date`
  createTunnel
fi


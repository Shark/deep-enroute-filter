#!/usr/bin/env bash
set -euo pipefail

GetWellKnown() {
  while :
  do
    coap get coap://[fdf0:a23f:8cae:5b97::2]/.well-known/core > /dev/null 2>&1
    sleep 1
  done
}

GetBasic() {
  while :
  do
    coap get coap://[fdf0:a23f:8cae:5b97::2]/basic > /dev/null 2>&1
    sleep 1
  done
}

main() {
  GetWellKnown &
  GetBasic &
}

main "$@"

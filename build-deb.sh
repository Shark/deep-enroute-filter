#!/usr/bin/env bash
set -euxo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

check_pkg() {
  declare name="$1"

  if dpkg --get-selections | grep "$name" > /dev/null; then
    return 0
  else
    return 1
  fi
}

build() {
  export PATH=$PATH:"$DIR/go/bin":"$DIR/gopath/bin" GOPATH="$DIR/gopath"

  if [[ ! -d "$DIR/go" ]]; then
    curl -sSL -o "$DIR/go.tar.gz" "https://dl.google.com/go/go1.11.linux-$ARCH.tar.gz"
    tar -C "$DIR/" -xf "$DIR/go.tar.gz"
    go get -u github.com/golang/dep/cmd/dep
  fi

  if [[ ! -d "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering" ]]; then
    mkdir -p "$GOPATH/src/gitlab.hpi.de/felix.seidel"
    cp -r "$DIR/iotsec-enroute-filtering" "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering"
  fi

  cd "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter"
  dep ensure
  go build .
}

package() {
  [[ -d "$DIR/package" ]] && rm -rf "$DIR/package"

  mkdir -p "$DIR/package/enroute-filter-0.1"
  cd "$DIR/package/enroute-filter-0.1"

  cp "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/filter" "$DIR/package/enroute-filter-0.1/enroute-filter"
  cp "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/enroute-filter.service" "$DIR/package/enroute-filter-0.1/"

  dh_make --createorig --packageclass s # single binary
  printf "enroute-filter usr/bin\nenroute-filter.service lib/systemd/system\n" > "$DIR/package/enroute-filter-0.1/debian/install"
  cat > "$DIR/package/enroute-filter-0.1/debian/postinst" <<EOF
#!/bin/bash
systemctl daemon-reload
systemctl enable enroute-filter
systemctl start enroute-filter
EOF
  chmod +x "$DIR/package/enroute-filter-0.1/debian/postinst"
  debuild -us -uc
  mv "$DIR"/package/*.deb "$DIR/"
}

main() {
    if [[ -z "${ARCH:-}" ]]; then
    >&2 echo "Please set ARCH to arm64, armv6l, or amd64"
    exit 1
  fi

  if [[ -z "${DEBEMAIL:-}" ]]; then
    >&2 echo "Please set DEBEMAIL"
    exit 1
  fi

  if [[ -z "${DEBFULLNAME:-}" ]]; then
    >&2 echo "Please set DEBFULLNAME"
    exit 1
  fi

  check_pkg 'libnetfilter-queue-dev'
  check_pkg 'git'
  check_pkg 'build-essential'
  check_pkg 'dh-make'
  check_pkg 'devscripts'

  if [[ ! -d "$DIR/iotsec-enroute-filtering" ]]; then
    >&2 echo "Please run 'git clone git@gitlab.hpi.de:felix.seidel/iotsec-enroute-filtering.git'"
    exit 1
  fi

  build
  package
}

main "$@"
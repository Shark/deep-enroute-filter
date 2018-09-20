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
  export PATH=$PATH:"$DIR/_build/go/bin":"$DIR/_build/gopath/bin" GOPATH="$DIR/_build/gopath"

  if [[ ! -d "$DIR/_build/go" ]]; then
    curl -sSL -o "$DIR/_build/go.tar.gz" "https://dl.google.com/go/go1.11.linux-$ARCH.tar.gz"
    tar -C "$DIR/_build/" -xf "$DIR/_build/go.tar.gz"
    go get -u github.com/golang/dep/cmd/dep
  fi

  if [[ ! -L "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering" ]]; then
    rm -rf "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering"
    mkdir -p "$GOPATH/src/gitlab.hpi.de/felix.seidel"
    ln -s "$DIR" "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering"
  fi

  cd "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter"
  dep ensure
  go build .
}

package() {
  [[ -d "$DIR/_build/package" ]] && rm -rf "$DIR/_build/package"

  mkdir -p "$DIR/_build/package/enroute-filter-$PKGVER"
  cd "$DIR/_build/package/enroute-filter-$PKGVER"

  cp "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/filter" "$DIR/_build/package/enroute-filter-$PKGVER/enroute-filter"
  cp "$GOPATH/src/gitlab.hpi.de/felix.seidel/iotsec-enroute-filtering/filter/enroute-filter.service" "$DIR/_build/package/enroute-filter-$PKGVER/"

  dh_make --yes --createorig --packageclass s # single binary
  printf "enroute-filter usr/bin\nenroute-filter.service lib/systemd/system\n" > "$DIR/_build/package/enroute-filter-$PKGVER/debian/install"
  cat > "$DIR/_build/package/enroute-filter-$PKGVER/debian/postinst" <<EOF
#!/usr/bin/env bash
systemctl daemon-reload
systemctl enable enroute-filter
systemctl start enroute-filter
EOF
  cat > "$DIR/_build/package/enroute-filter-$PKGVER/debian/prerm" <<EOF
#!/usr/bin/env bash
systemctl stop enroute-filter
systemctl disable enroute-filter
EOF
  cat > "$DIR/_build/package/enroute-filter-$PKGVER/debian/postrm" <<EOF
#!/usr/bin/env bash
systemctl daemon-reload
EOF
  chmod +x "$DIR/_build/package/enroute-filter-$PKGVER/debian"/{postinst,prerm,postrm}
  debuild -us -uc
  mv "$DIR"/_build/package/*.deb "$DIR/_build/"
}

main() {
  if [[ -z "${ARCH:-}" ]]; then
    printf "ARCH [arm64, armv6l, amd64]: "
    read ARCH
    if [[ -z $ARCH ]]; then
      >&2 echo "No ARCH given, exiting"
      exit 1
    fi
    export ARCH
  fi

  if [[ -z "${DEBEMAIL:-}" ]]; then
    printf "DEBEMAIL [some@email.com]: "
    read DEBEMAIL
    if [[ -z $DEBEMAIL ]]; then
      >&2 echo "No DEBEMAIL given, exiting"
      exit 1
    fi
    export DEBEMAIL
  fi

  if [[ -z "${DEBFULLNAME:-}" ]]; then
    printf "DEBFULLNAME [Some Name]: "
    read DEBFULLNAME
    if [[ -z $DEBFULLNAME ]]; then
      >&2 echo "No DEBFULLNAME given, exiting"
      exit 1
    fi
    export DEBFULLNAME
  fi

  if [[ -z "${PKGVER:-}" ]]; then
    printf "PKGVER [0.1.0-1]: "
    read PKGVER
    if [[ -z $PKGVER ]]; then
      >&2 echo "No PKGVER given, exiting"
      exit 1
    fi
    export PKGVER
  fi

  check_pkg 'libnetfilter-queue-dev'
  check_pkg 'git'
  check_pkg 'build-essential'
  check_pkg 'dh-make'
  check_pkg 'devscripts'

  rm -rf "$DIR/_build/package" "$DIR"/_build/*.deb
  build
  package
}

main "$@"
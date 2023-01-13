#!/bin/sh

uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  echo "$os"
}


uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    x86) arch="386" ;;
    i686) arch="386" ;;
    i386) arch="386" ;;
    aarch64) arch="arm64" ;;
    armv*) arch="armv6" ;;
    armv*) arch="armv6" ;;
    armv*) arch="armv6" ;;
  esac
  echo ${arch}
}

TAG="v0.0.1"
OS=$(uname_os)
ARCH=$(uname_arch)
INSTALL_DIR="/usr/local/bin/"

BINARY_URL="https://github.com/shebogholo/shex/releases/download/$TAG/shex-$ARCH-$OS"

tmp=$(mktemp -d)

binary=$tmp/shex

echo "===> Downloading..."
$(curl $BINARY_URL -sL -o $binary)

echo "===> Installing..."
sudo install "$binary" $INSTALL_DIR

rm -rf "$tmp"

echo "Installed successfully!"
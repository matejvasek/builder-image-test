#!/usr/bin/env sh

set -ex

mkdir -p "$HOME/bin/"

# install latest from main branch
TMP_DIR="$(mktemp -d)"
cd "$TMP_DIR"
go get github.com/markbates/pkger/cmd/pkger
git clone https://github.com/boson-project/func
cd func
make
cp func "$HOME/bin/func_latest"
cd "$HOME"
rm -fr "$TMP_DIR"

# install latest official release
curl -L -o - https://github.com/boson-project/func/releases/latest/download/func_linux_amd64.gz | gunzip > "$HOME/bin/func_stable"
chmod 755 "$HOME/bin/func_stable"

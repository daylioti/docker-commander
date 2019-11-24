#/bin/bash

cmd="curl -sL -o"


if [ "${_GOOS}" == "linux" ] ||  [ "${_GOOS}" == "darwin" ]; then
  cmd="${cmd} upx.tar.xz https://github.com/upx/upx/releases/download/v3.95/upx-3.95-"
  cmd="${cmd}${_GOARCH}_linux.tar.xz"
  $cmd
  unxz upx.tar.xz
  tar -xvf upx.tar
  cp upx-3.95*/upx .
fi

if [ "${_GOOS}" == "windows" ]; then
  cmd="${cmd} upx.zip https://github.com/upx/upx/releases/download/v3.95/upx-3.95-"
  cmd="${cmd}win"
  if [ "$_GOARCH" == "386" ]; then
    cmd="${cmd}32.zip"
  fi
  if [ "$_GOARCH" == "amd64" ]; then
    cmd="${cmd}64.zip"
  fi
  $cmd
  unzip upx.zip
  cp upx*/upx .
fi

ls -lah
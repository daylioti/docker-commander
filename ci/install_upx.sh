#/bin/bash

cmd="curl -sL -o"


if [ "${_GOOS}" == "linux" ] ||  [ "${_GOOS}" == "darwin" ]; then
  cmd="${cmd} upx.tar.xz https://github.com/upx/upx/releases/download/v3.95/upx-3.95-"
  if [ "${_GOARCH}" == "386" ]; then
    cmd="${cmd}i${_GOARCH}_linux.tar.xz"
  else
    cmd="${cmd}${_GOARCH}_linux.tar.xz"
  fi
  $cmd
  unxz upx.tar.xz
  tar -xvf upx.tar
  cp upx-3.95*/upx .
fi

#/bin/bash

cmd="curl -sL -o"


if [ "${_GOOS}" == "linux" ]; then
  cmd="${cmd} upx.tar.xz https://github.com/upx/upx/releases/download/v3.95/upx-3.95-"
  cmd="${cmd}${_GOARCH}_linux.tar.xz"
  $cmd
  unxz upx.tar.xz
  tar -xvf upx.tar
  cp upx-3.95*/upx .
fi


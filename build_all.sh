#!/usr/bin/env bash



CURDIR=`/bin/pwd`
BASEDIR=$(dirname $0)
ABSPATH=$(readlink -f $0)
ABSDIR=$(dirname $ABSPATH)

cd $ABSDIR/../../../../
GOPATH=`pwd`
cd $CURDIR
bash $ABSDIR/build_package.sh "github.com/arnaucode/derosuite/cmd/derod"
bash $ABSDIR/build_package.sh "github.com/arnaucode/derosuite/cmd/explorer"
bash $ABSDIR/build_package.sh "github.com/arnaucode/derosuite/cmd/dero-wallet-cli"
cd "${ABSDIR}/build"

#all other platforms are okay with tar.gz
find . -mindepth 1 -type d -not -name '*windows*'   -exec tar --owner=dummy --group=dummy -cvzf {}.tar.gz {} \;

cd $CURDIR

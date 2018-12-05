#!/bin/bash

MYOXI_DIR='./.myoxi-release'
VERSION=`git describe --long --tags --dirty --always | sed -e 's/^v//'`
for os in linux windows
do
    for arch in amd64 386
    do
        rm -Rf ${MYOXI_DIR}
        GOOS=$os GOARCH=$arch go build -ldflags "-X main.MyoxiVersion=$VERSION" .

        NAME=myoxi-${VERSION}-${os}-${arch}
        REL_DIR=${MYOXI_DIR}/${NAME}
        mkdir -p ${REL_DIR}
        cp ./myoxi* ${REL_DIR}/
        cp ./README.md ${REL_DIR}/
        cp ./CHANGELOG.md ${REL_DIR}/
        cp ./LICENSE ${REL_DIR}/
        cd ${MYOXI_DIR} && zip -r ../${NAME}.zip ${NAME}
        cd ..
        rm -Rf ${MYOXI_DIR}
        rm -f myoxi myoxi.exe
    done
done

#!/bin/bash

export LIBGADU_VERSION="1.12.1"

if [[ $TRAVIS_OS_NAME == 'osx' ]]; then
    # Install some custom requirements on OS X
    brew update
    brew install libgadu
else
    # Install some custom requirements on Linux
    pushd $PWD
    cd /tmp
    wget http://github.com/wojtekka/libgadu/releases/download/$LIBGADU_VERSION/libgadu-$LIBGADU_VERSION.tar.gz
    tar -zxvf libgadu-$LIBGADU_VERSION.tar.gz
    cd libgadu-$LIBGADU_VERSION
    ./configure --prefix=/usr
    make
    sudo make install
    popd
fi

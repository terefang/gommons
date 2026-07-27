#!/usr/bin/env just --justfile

XDIR := justfile_directory()
XDIST := XDIR + "/dist"
EXE := "gommons"
CLI := "cmdbox"
OSARCH := os()+"-"+arch()

GOARCHS := "linux,amd64 linux,arm64 windows,amd64 windows,arm64 darwin,amd64 darwin,arm64"

update-version-info:
    #!/bin/sh
    V=$(shtool version -n "{{EXE}}" -d long -l txt ./version.txt)
    echo "$V" >{{XDIR}}/version_info.txt
    echo "package pkg" >{{XDIR}}/pkg/version.go
    echo "const (" >>{{XDIR}}/pkg/version.go
    echo "    PkgVersion = \"$V\"" >>{{XDIR}}/pkg/version.go
    echo ")" >>{{XDIR}}/pkg/version.go

set-drel: do-inc-level
    #!/bin/bash
    V=$(date '+%Y.%m.')
    V=$V$(cd {{XDIR}} && shtool version -n "{{EXE}}" -l short ./version.txt|cut -f 3 -d.)
    cd {{XDIR}} && just -f justfile set-version "$V"

git-push-level: inc-level git-push

git-push: (git-push-msg "upd")
git-push-msg _MSG:
    #!/bin/sh
    VERSION=$(shtool version -l txt ./version.txt)
    TIME=$(date '+%F %T')
    git commit --all -m "{{_MSG}} - v$VERSION on $TIME"
    git push

inc-version: do-inc-version update-version-info
do-inc-version:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i v -l txt ./version.txt

inc-major: do-inc-major update-version-info
do-inc-major:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i r -l txt ./version.txt

inc-level: do-inc-level update-version-info
do-inc-level:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i l -l txt ./version.txt

set-version _VERSION: (do-set-version _VERSION) update-version-info

do-set-version _VERSION:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -s "{{_VERSION}}" -l txt ./version.txt

make-update: git-push-level

make-prel: (make-prel-msg "upd")
make-prel-msg _MSG: set-drel git-push
    #!/bin/bash
    VERSION=$(shtool version -l txt ./version.txt)
    VERL=$(shtool version -l text -d long ./version.txt)
    MESSAGE="{{_MSG}} - {{EXE}} automated pre-release version $VERL"
    gh release create v$VERSION --notes "$MESSAGE" --prerelease

make-rel: (make-rel-msg "upd")
make-rel-msg _MSG: set-drel git-push
    #!/bin/bash
    VERSION=$(shtool version -l txt ./version.txt)
    VERL=$(shtool version -l text -d long ./version.txt)
    MESSAGE="{{_MSG}} - {{EXE}} automated release version $VERL"
    gh release create v$VERSION --notes "$MESSAGE"

make-upload:
    #!/bin/bash
    VERSION=$(shtool version -l txt ./version.txt)
    gh release upload v$VERSION {{XDIST}}/*

make-prel-full _MSG: (make-prel-msg _MSG) build make-upload
make-rel-full _MSG: (make-rel-msg _MSG) build make-upload

build: update-version-info
    #!/bin/sh
    export GOROOT=${HOME}/bin/go
    export PATH=$GOROOT/bin:$PATH
    VERSION=$(shtool version -l txt ./version.txt)
    rm -rf {{XDIST}}; mkdir -p {{XDIST}}
    #CGO_ENABLED=1 CC=musl-gcc \
    #  go build -ldflags="-linkmode external -extldflags '-static'"
    for XTYPE in {{GOARCHS}}; do
        _gos=${XTYPE%%,*}
        _garch=${XTYPE##*,}
        _gexe=
        case $_gos in
            windows)
                _gexe=.exe
                ;;
            *)
                _gexe=
                ;;
        esac
        echo "building ${_gos} ${_garch} ..."
        GOOS=${_gos} GOARCH=${_garch} CGO_ENABLED=0 go build -buildmode=pie -o "{{XDIST}}/{{CLI}}-${VERSION}-${_gos}-${_garch}${_gexe}" cmd/{{CLI}}/*.go
    done

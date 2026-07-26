#!/usr/bin/env just --justfile

XDIR := justfile_directory()
EXE := "gommons"
CLI := "cmdbox"
OSARCH := os()+"-"+arch()

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

git-push:
    #!/bin/sh
    VERSION=$(shtool version -l txt ./version.txt)
    TIME=$(date '+%F %T')
    git commit --all -m "upd to v$VERSION on $TIME"
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

make-prel: set-drel git-push
    #!/bin/bash
    VERSION=$(shtool version -l txt ./version.txt)
    VERL=$(shtool version -l text -d long ./version.txt)
    MESSAGE="{{EXE}} automated pre-release version $VERL"
    gh release create v$VERSION --notes "$MESSAGE" --prerelease out/*

make-rel: set-drel git-push
    #!/bin/bash
    VERSION=$(shtool version -l txt ./version.txt)
    VERL=$(shtool version -l text -d long ./version.txt)
    MESSAGE="{{EXE}} automated release version $VERL"
    gh release create v$VERSION --notes "$MESSAGE" out/*

build: update-version-info
    #!/bin/sh
    export GOROOT=${HOME}/bin/go
    export PATH=$GOROOT/bin:$PATH
    VERSION=$(shtool version -l txt ./version.txt)
    rm -rf {{XDIR}}/out; mkdir -p {{XDIR}}/out
    go build -o "{{XDIR}}/out/{{CLI}}-${VERSION}-{{OSARCH}}" cmd/{{CLI}}/*.go
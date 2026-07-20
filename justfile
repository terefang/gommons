#!/usr/bin/env just --justfile

XDIR := justfile_directory()
EXE := "gommons"

set-drel: inc-level
    #!/bin/bash
    V=$(date '+%Y.%m.')
    V=$V$(cd {{XDIR}} && shtool version -n "{{EXE}}" -l short ./version.txt|cut -f 3 -d.)
    cd {{XDIR}} && just -f justfile set-version "$V"

inc-version:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i v -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

inc-major:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i r -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

inc-level:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -i l -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

set-version _VERSION:
    #!/bin/bash
    cd {{XDIR}}
    shtool version -n "{{EXE}}" -s "{{_VERSION}}" -l txt ./version.txt
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt

make-release: set-drel build
    #!/bin/bash
    VERSION=$(shtool version -l txt ./version.txt)
    MESSAGE="automated release version $(shtool version -l text -d long ./version.txt)"
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt
    gh release create v$VERSION --notes "$MESSAGE"

build:
    #!/bin/sh
    export GOROOT=${HOME}/bin/go
    export PATH=$GOROOT/bin:$PATH
    VERSION=$(shtool version -l txt ./version.txt)
    shtool version -n "{{EXE}}" -d long -l txt ./version.txt >{{XDIR}}/version_info.txt
    #go build
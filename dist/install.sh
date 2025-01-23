#!/bin/sh
cd "$(dirname -- "$0")" || exit 1

install -vDm0755 "bin/hizla" "${HIZLA_INSTALL_PREFIX}/usr/bin/hizla"

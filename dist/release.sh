#!/bin/sh -e
cd "$(dirname -- "$0")/.."
VERSION="${HIZLA_VERSION:-untagged}"
pname="hizla-${VERSION}"
out="dist/${pname}"

mkdir -p "${out}"
cp -v "dist/install.sh" "${out}"

go generate ./...
go build -trimpath -v -o "${out}/bin/" -ldflags "-s -w -buildid=
  -X github.com/hizla/hizla/internal.Version=${VERSION}" ./...

rm -f "./${out}.tar.gz" && tar -C dist -czf "${out}.tar.gz" "${pname}"
rm -rf "./${out}"
(cd dist && sha512sum "${pname}.tar.gz" >"${pname}.tar.gz.sha512")

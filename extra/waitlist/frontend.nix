{
  lib,
  stdenvNoCC,
  fetchFromGitHub,
  mkYarnPackage,
  yarn,
  ibm-plex,
  inter,

  baseUrl ? "http://localhost:3000/",
  # https://github.com/vercel/next.js/issues/10111
  cacheDir ? "/tmp/nextjs",
}:

stdenvNoCC.mkDerivation rec {
  pname = "hizla-frontend";
  version = "0.0.2";

  src = fetchFromGitHub {
    owner = "hizla";
    repo = "waitlist-frontend";
    rev = "refs/tags/v${version}";
    sha256 = "sha256-uI/hpVAeMnGPBDt1B0F+2W/ekZugbHuoFxUkwUmf99U=";
  };

  patches = [ ./local-fonts.patch ];

  nodeModules = mkYarnPackage {
    pname = "${pname}-node-modules";
    inherit version src;
  };

  NEXT_PUBLIC_API = baseUrl;

  buildPhase = ''
    ln -s ${ibm-plex}/share/fonts/opentype ibm-plex
    ln -s ${inter}/share/fonts/truetype/InterVariable.ttf

    ln -s ${nodeModules}/libexec/${pname}/node_modules

    HOME="$(mktemp -d)" ${lib.getExe yarn} --offline build
  '';

  installPhase = ''
    mkdir -p $out/share
    cp -r ${src} $out/share/${pname}
    chmod +w $out/share/${pname}
    cp -r ./.next $out/share/${pname}/.next
    rm -rf $out/share/${pname}/.next/cache
    ln -s ${cacheDir} $out/share/${pname}/.next/cache
    ln -s ${nodeModules}/libexec/${pname}/node_modules $out/share/${pname}/

    mkdir -p $out/bin
    exe="$out/bin/${pname}"
    echo "#!/usr/bin/env bash" > $exe
    echo "cd $out/share/${pname}" >> $exe
    echo "${lib.getExe yarn} run start" >> $exe
    chmod +x $exe
  '';
}

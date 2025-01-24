{
  symlinkJoin,

  backend,
  frontend,

  # frontend workaround values
  baseUrl ? "",
  cacheDir ? "/tmp/nextjs",
}:

symlinkJoin {
  name = "hizla-waitlist";
  paths = [
    backend
    (frontend { inherit baseUrl cacheDir; })
  ];

  meta.mainProgram = backend.pname;
}

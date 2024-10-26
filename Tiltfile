# Install Node deps on change of package.json
local_resource(
  'npm',
  cmd='npm i',
  deps=['package.json'],
)

# Rebuild frontend if source files change
local_resource(
  'frontend',
  cmd='node ./ci/build.mjs',
  deps=['src'],
  resource_deps=['npm'],
)

# Generate translation files on source change
local_resource(
  'translations',
  cmd='make translate',
  deps=['i18n.yaml'],
)

# Rebuild and run Go webserver on code changes
local_resource(
  'server',
  cmd='go build .',
  deps=[
    'api.go',
    'frontend',
    'helpers.go',
    'main.go',
    'pkg',
    'storage.go',
    'tplFuncs.go',
  ],
  ignore=['ots', 'src'],
  serve_cmd='./ots --listen=:15641',
  serve_env={
    'CUSTOMIZE': 'customize.yaml',
  },
  readiness_probe=probe(
    http_get=http_get_action(15641, path='/api/healthz'),
    initial_delay_secs=1,
  ),
  resource_deps=['frontend', 'translations'],
)

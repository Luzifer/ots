# Install Node deps on change of package.json
local_resource(
  'yarn',
  cmd='corepack yarn@1 install', # Not using the make target to edit the lockfile
  deps=['package.json'],
)

# Rebuild frontend if source files change
local_resource(
  'frontend',
  cmd='make frontend',
  deps=['src'],
  resource_deps=['yarn'],
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
  deps=[
    'api.go',
    'frontend',
    'helpers.go',
    'main.go',
    'pkg',
    'storage.go',
    'tplFuncs.go',
    'go.mod', 'go.sum',
  ],
  ignore=[
    'src'
  ],
  serve_cmd='go run . --listen=:15641',
  serve_env={
    'CUSTOMIZE': 'customize.yaml',
  },
  readiness_probe=probe(
    http_get=http_get_action(15641, path='/api/healthz'),
    initial_delay_secs=1,
  ),
  resource_deps=[
    'frontend',
    'translations',
  ],
)

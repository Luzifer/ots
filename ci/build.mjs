import esbuild from 'esbuild'
import { sassPlugin } from 'esbuild-sass-plugin'
import vuePlugin from 'esbuild-plugin-vue3'

esbuild.build({
  assetNames: '[name]',
  bundle: true,
  define: {
    'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'dev'),
  },
  entryPoints: ['src/main.ts'],
  legalComments: 'none',
  loader: {
    '.ttf': 'empty', // Drop files, all targets below support woff2
    '.woff2': 'file',
  },
  minify: true,
  outfile: 'frontend/app.js',
  plugins: [
    sassPlugin(),
    vuePlugin(),
  ],
  target: [
    'chrome109',
    'edge132',
    'es2020',
    'firefox115',
    'safari16',
  ],
})

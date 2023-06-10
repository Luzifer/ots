import { sassPlugin } from 'esbuild-sass-plugin'
import vuePlugin from 'esbuild-vue'
import esbuild from 'esbuild'

esbuild.build({
  assetNames: '[name]',
  bundle: true,
  define: {
    'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'dev'),
  },
  entryPoints: ['src/main.js'],
  loader: {
    '.woff2': 'file',
  },
  minify: true,
  outfile: 'frontend/app.js',
  plugins: [
    sassPlugin(),
    vuePlugin(),
  ],
  target: [
    'chrome87',
    'edge87',
    'es2020',
    'firefox84',
    'safari14',
  ],
})

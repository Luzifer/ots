const path = require('path')
const webpack = require('webpack');

const VueLoaderPlugin = require('vue-loader/lib/plugin')

function resolve(dir) {
  return path.join(__dirname, dir)
}

module.exports = {
  entry: './main.js',
  output: {
    filename: 'app.js',
    path: path.resolve(__dirname, '..', 'frontend')
  },
  plugins: [
    new webpack.optimize.OccurrenceOrderPlugin(),
    new webpack.DefinePlugin({
      'process.env': {
        'NODE_ENV': JSON.stringify('production')
      }
    }),
    new VueLoaderPlugin(),
  ],
  optimization: {
    minimize: true
  },
  module: {
    rules: [

      {
        test: /\.(s?)css$/,
        use: [
          'style-loader',
          'css-loader',
          'sass-loader',
        ],
      },

      {
        test: /\.js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: [
              ['env', { "targets": { "browsers": [">0.25%", "not ie 11", "not op_mini all"] } }]
            ]
          }
        }
      },

      {
        test: /\.vue$/,
        loader: 'vue-loader',
      },

    ]
  }
}

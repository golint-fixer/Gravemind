const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
  resolve: {
    root: path.resolve(__dirname, 'frontend'),
    extensions: ['', '.js', '.jsx', '.json'],
  },
  entry: [
    './frontend/index.jsx',
  ],
  output: {
    path: path.resolve(__dirname, 'build'),
    filename: 'bundle.js',
    publicPath: '/build/',
  },
  plugins: [
    new webpack.NoErrorsPlugin(),
    new HtmlWebpackPlugin({
      template: 'frontend/index.jade',
    }),
  ],
  module: {
    loaders: [
      {
        test: /static/,
        loader: 'file?name=[path][name].[ext]&context=frontend/static',
      },
      {
        test: /\.(png|jpg|jpeg|gif|svg|woff|woff2)$/,
        loader: 'url-loader?limit=10000',
        exclude: /static/,
      },
      {
        test: /\.scss$/,
        loaders: [
          'style-loader',
          'css-loader',
          'postcss-loader?parser=postcss-scss',
        ],
      },
      {
        test: /\.jade$/,
        loader: 'jade',
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: 'babel',
      },
      {
        test: /\.jsx$/,
        exclude: /node_modules/,
        loader: 'babel',
        query: {
          plugins: [
            [
              'react-transform',
              {
                'transforms': [
                  {
                    'transform': 'react-transform-hmr',
                    'imports': ['react'],
                    'locals': ['module'],
                  },
                  {
                    'transform': 'react-transform-catch-errors',
                    'imports': ['react', 'redbox-react'],
                  },
                ],
              },
            ],
          ],
        },
      },
    ],
  },
  postcss: (bundler) => {
    return [
      require('postcss-import')({ addDependencyTo: bundler }),
      require('precss')(),
      require('autoprefixer')(),
    ];
  },
};

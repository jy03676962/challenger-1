var path = require('path');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin')
var CleanWebpackPlugin = require('clean-webpack-plugin')

module.exports = {
  devtool: 'eval',
  entry: [
    // 'webpack-dev-server/client?http://localhost:3000',
    // 'webpack/hot/only-dev-server',
    './src/app.jsx'
  ],
  output: {
    path: path.join(__dirname, 'dist/js'),
    filename: 'bundle.js',
    publicPath: '/js/'
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new HtmlWebpackPlugin({
      template: 'src/assets/index.ejs',
      inject: 'body',
      filename: '../index.html'
    }),
    new HtmlWebpackPlugin({
      template: 'src/assets/api.ejs',
      inject: 'body',
      filename: '../api.html'
    }),
    new CleanWebpackPlugin(['dist'])
  ],
  module: {
    loaders: [{
      test: /\.jsx$/,
      loaders: ['react-hot', 'babel'],
      include: path.join(__dirname, 'src')
    }, {
      test: /\.css$/,
      loaders: [
        'style?sourceMap',
        'css?modules&importLoaders=1&localIdentName=[path]___[name]__[local]___[hash:base64:5]'
      ],
      include: path.join(__dirname, 'src/styles')
    }, {
      test: /\.(png|jpg|gif)$/,
      loader: "file-loader?name=../img/img-[hash:6].[ext]"
    }]
  }
};

const path = require('path');
const {CleanWebpackPlugin} = require('clean-webpack-plugin');
const MiniCSSExtractPlugin = require('mini-css-extract-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const TerserJSPlugin = require('terser-webpack-plugin');

module.exports = {
	optimization: {
		minimizer: [
			new TerserJSPlugin({
				extractComments: false
			}),
			new OptimizeCSSAssetsPlugin({})
		]
	},
	entry: {
		edit: './src/js/edit.js',
		config: './src/js/config.js',
		list: './src/js/list.js',
		signin: './src/js/signin.js'
	},
	output: {
		filename: '[name].bundle.js',
		path: path.resolve(__dirname,'dist')
	},
	module: {
		rules: [
			{
				test: /\.css$/i,
				use: [
					MiniCSSExtractPlugin.loader,
					'css-loader'
				]
			},{
				test: /\.png$/,
				loader: 'file-loader'
			},{
				test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
				loader: "url-loader?limit=10000&mimetype=application/font-woff"
			},{
				test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
				loader: "file-loader"
			}
		]
	},
	plugins: [
		new CleanWebpackPlugin({}),
		new MiniCSSExtractPlugin({
			filename: '[name].bundle.css'
		})
	]
};

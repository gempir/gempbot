const path = require("path");
const webpack = require("webpack");

module.exports = (env, options) => ({
	entry: './src/index.jsx',
	output: {
		path: path.resolve(__dirname, 'public'),
		filename: 'bundle.js',
		publicPath: "/",
	},
	module: {
		rules: [
			{
				test: /\.jsx$/,
				exclude: /node_modules/,
				use: {
					loader: "babel-loader"
				}
			},
			{
				test: /\.s?css$/,
				use: ["style-loader", "css-loader", "sass-loader"]
			},
		]
	},
	stats: {
		// Config for minimal console.log mess.
		assets: false,
		colors: true,
		version: false,
		hash: false,
		modules: false,
		timings: false,
		entrypoints: false,
		chunks: false,
		chunkModules: false
	},
	plugins: [
		new webpack.DefinePlugin({
			'process.env': {
				'apiBaseUrl': options.mode === 'development' ? '"http://localhost:8035"' : '"https://spampchamp-api.gempir.com"',
			}
		}),
	],
	resolve: {
		extensions: ['.js', '.jsx'],
	}
});
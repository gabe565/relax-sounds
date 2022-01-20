module.exports = {
  chainWebpack(config) {
    config.plugin('html').tap((args) => {
      args[0].title = 'Relax Sounds';
      return args;
    });
    config.module
      .rule('data')
      .test(/data\/.+\.json$/)
      .type('javascript/auto')
      .use('file-loader')
      .loader('file-loader')
      .options({
        name: '[path][name].[ext]',
        context: 'src',
      })
      .end();
  },
  pwa: {
    name: 'Relax Sounds',
    themeColor: '#673AB7',
    msTileColor: '#673AB7',
    appleMobileWebAppCapable: 'yes',
    appleMobileWebAppStatusBarStyle: 'black',

    workboxPluginMode: 'InjectManifest',
    workboxOptions: {
      swSrc: './src/service-worker.js',
      exclude: [
        /^data\/.*\.json$/,
        /^data\/audio\//,
        /\.map$/,
        /img\/icons\//,
        /favicon\.ico$/,
        /^manifest.*\.js?$/,
      ],
    },
  },
  transpileDependencies: [
    'vuetify',
  ],
};

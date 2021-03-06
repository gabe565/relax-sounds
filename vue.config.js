module.exports = {
  chainWebpack: (config) => {
    config.plugin('html').tap((args) => {
      args[0].title = 'Relax Sounds';
      return args;
    });
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
        /^audio\//,
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

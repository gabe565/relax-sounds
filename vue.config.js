module.exports = {
  chainWebpack: (config) => {
    config.plugin('html').tap((args) => {
      args[0].title = 'Relax Sounds';
      return args;
    });
  },
  pwa: {
    name: 'Relax Sounds',
    workboxOptions: {
      swSrc: './src/service-worker.js',
      swDest: 'service-worker.js',
    },
    workboxPluginMode: 'InjectManifest',
  },
  transpileDependencies: [
    'vuetify',
  ],
};

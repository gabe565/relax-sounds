module.exports = {
  root: true,
  env: {
    browser: true,
  },
  extends: [
    'airbnb-base',
    'plugin:vue/vue3-recommended',
    'plugin:vuetify/base',
    'prettier',
    'plugin:prettier/recommended',
  ],
  parserOptions: {
    sourceType: 'module',
  },
  rules: {
    'import/prefer-default-export': 'off',
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-param-reassign': 'off',
    'no-underscore-dangle': 'off',
    'no-restricted-syntax': ['error', 'ForInStatement', 'LabeledStatement', 'WithStatement'],
    'import/no-unresolved': 'off',
    "import/extensions": "off",
  },
};

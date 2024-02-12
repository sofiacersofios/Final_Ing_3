exports.config = {
  tests: './*_test.js',
  output: './output',
  helpers: {
    Puppeteer: {
      url: 'http://localhost:3000',
      show: true,
      windowSize: '1200x900',
      chrome: {
        args: [
          // Agrega la opción --disable-web-security para deshabilitar CORS
          '--disable-web-security',
        ],
      },
    },
  },
  include: {
    I: './steps_file.js',
  },
  bootstrap: null,
  teardown: null,
  mocha: {},
  name: 'frontend',
  plugins: {
    screenshotOnFail: {
      enabled: true,
    },
  },
};

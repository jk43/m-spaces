module.exports = function (ctx) {
  return {
    devServer: {
      https: {
        key: 'certs/key.pem',
        cert: 'certs/cert.pem',
        ca: 'certs/chain.pem',
      },
    },
  }
}

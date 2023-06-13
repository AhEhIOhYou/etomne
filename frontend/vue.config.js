const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    proxy: {
      '/api': {
        target: 'http://model3d-api:8095',
        secure: false,
        changeOrigin: true,
      },
    }
  }
})

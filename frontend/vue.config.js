const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    proxy: {
      '/api': {
        target: 'https://modelshowtime.serdcebolit.ru',
        secure: false,
        changeOrigin: true,
      },
    }
  }
})

<template>
  <header class="header">
    <navbar></navbar>
  </header>
  <div class="container">
    <router-view/>
  </div>
</template>

<script>
import axios from "axios";
import VueCookies from 'vue-cookies';

export default {  
  watch: {
    $route: 'fetchData',
  },
  methods: {
    fetchData() {
      const accessToken = $cookies.isKey('access_token');
      const refreshToken = $cookies.isKey('refresh_token');
      if (!accessToken && refreshToken) {
        axios.post('/api/users/refresh', {
            refresh_token: refreshToken
          })
        .then(response => {
          $cookies.set('access_token', response.data.access_token, '15m', '/');
          $cookies.set('refresh_token', response.data.refresh_token, '7d', '/');
          commit('setId', `${response.data.id}`);
          commit('setName', `${response.data.name}`);
          localStorage.setItem('name', response.data.name);
          localStorage.setItem('id', response.data.id);
        })
        .catch(error => {
          console.log(error);
        });
      }
    }
  }
}
</script>

<style lang="scss">
@import "@/assets/styles/vendor/_normalize.scss";
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
@import "@/assets/styles/global/_fonts.scss";
@import "@/assets/styles/global/_reboot.scss";
@import "@/assets/styles/global/_utils.scss";
@import "@/assets/styles/global/_container.scss";
@import "@/assets/styles/blocks/_header.scss";
@import "@/assets/styles/blocks/_btn.scss";
@import "@/assets/styles/blocks/_model.scss";
</style>

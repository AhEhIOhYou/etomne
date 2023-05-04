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
  data () {
    return {
    }
  },
  watch: {
    $route: 'fetchData',
  },
  methods: {
    fetchData() {
      const isAccessTokenExist = $cookies.isKey('access_token');
      const isRefreshTokenExit = $cookies.isKey('refresh_token');
      if (!isAccessTokenExist && isRefreshTokenExit) {
        const refreshToken = $cookies.get('refresh_token');
        axios.post('/api/users/refresh', {
            refresh_token: refreshToken
          })
        .then(response => {
          $cookies.set('access_token', response.data.tokens.access_token, '15min', '/');
          $cookies.set('refresh_token', response.data.tokens.refresh_token, '7d', '/');
          commit('setId', `${response.data.public_data.id}`);
          commit('setName', `${response.data.public_data.name}`);
          localStorage.setItem('name', response.data.public_data.name);
          localStorage.setItem('id', response.data.public_data.id);
          localStorage.setItem('isAuth', false);
        })
        .catch(error => {
          console.log(error);
        });
      } else {
        if (localStorage.name) {
          localStorage.setItem('isAuth', true);
        } else {
          localStorage.setItem('isAuth', false);
        }
      }
    }
  },
  created() {
    this.fetchData();
  },
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

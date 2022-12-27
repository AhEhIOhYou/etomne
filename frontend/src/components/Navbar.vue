<template>
  <nav class="main-nav">
    <div class="main-nav__wrapper">
      <ul class="main-nav__list">
        <router-link class="main-nav__logo" to="/"><img src="@/assets/logo.png" width="640" height="640" alt="Лого"/></router-link>
        <li class="main-nav__item">
          <router-link v-show="isAuth" class="main-nav__link btn btn--white" to="/uploadmodel">Загрузка моделей</router-link>
        </li>
        <li class="main-nav__item">
          <router-link v-if="!isAuth" class="main-nav__link btn" to="/authorization">Авторизация</router-link>
          <button v-else @click="logout" class="main-nav__link btn" type="button">Выйти</button>
        </li>
      </ul>
    </div>
  </nav>
</template>
<script>
import axios from "axios";
// import VueCookies from 'vue-cookies';
// import auth from '@/assets/functions/auth.js';

export default {  
  name: 'navbar',
  props: {
    isAuth: Boolean,
  },
  methods: {
    logout() {
      const accessToken = $cookies.get("access_token");
      const refreshToken = $cookies.get("refresh_token");

      const logoutFunc = (access) => {
        axios.get("/api/users/logout/", {
          headers: {
            "Authorization": `Bearer ${access}`
          }
        })
          .then(res => {
            $cookies.remove("access_token");
            $cookies.remove("refresh_token");
            localStorage.removeItem('name');
            localStorage.removeItem('id');
            window.location.href = '/authorization';
          })
          .catch(error => {
            console.log(error)
          });
      };

      if (accessToken === null && refreshToken) {
        Promise.all([axios.post('/api/users/refresh', {
          refresh_token: refreshToken
        })
          .then(response => {
            $cookies.set('access_token', response.data.access_token, '15m', '/');
            $cookies.set('refresh_token', response.data.refresh_token, '7d', '/');
            commit('setId', `${response.data.id}`);
            commit('setName', `${response.data.name}`);
            localStorage.setItem('name', response.data.name);
            localStorage.setItem('id', response.data.id);
            return { newAccessToken: response.data.access_token, newRefreshToken: response.data.refresh_token }
          })
          .catch(error => {
            console.log(error);
          })]).then(function (values) {
            logoutFunc(values[0].newAccessToken);
          });
      } else {
        logoutFunc(accessToken);
      }
    }
  }
}
</script>
<style scoped lang="scss">
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
@import "@/assets/styles/blocks/_main-nav.scss";
</style>
<template>
  <nav class="main-nav">
    <div class="main-nav__wrapper">
      <ul class="main-nav__list">
        <router-link class="main-nav__logo" to="/"><img src="@/assets/logo.png" width="50" height="50" alt="Лого"/></router-link>
        <li class="main-nav__item">
          <router-link v-show="isAuth" class="main-nav__link btn btn--white" to="/uploadmodel">Загрузка моделей</router-link>
        </li>
        <li class="main-nav__item">
          <router-link v-if="!isAuth" class="main-nav__link btn" to="/authorization">Авторизация</router-link>
          <router-link v-if="isAuth" class="main-nav__link btn" to="/lk">Личный кабинет</router-link>
          <button v-if="isAuth" @click="logout" class="main-nav__link btn btn--white" type="button">Выйти</button>
        </li>
      </ul>
    </div>
  </nav>
</template>
<script>
import axios from "axios";
import VueCookies from 'vue-cookies';

export default {  
  name: 'navbar',
  props: {
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
            localStorage.setItem('isAuth', false);
            window.location.href = '/authorization';
          })
          .catch(error => {
            console.log(error)
          });
      };

      if (accessToken === null && refreshToken) {
        axios.post('/api/users/refresh', {
          refresh_token: refreshToken
        })
          .then(response => {
            $cookies.set('access_token', response.data.tokens.access_token, '15min', '/');
            $cookies.set('refresh_token', response.data.tokens.refresh_token, '7d', '/');
            localStorage.setItem('isAuth', true);
            logoutFunc(response.data.tokens.access_token);
            console.log(response);
          })
          .catch(error => {
            console.log(error);
          });
      } else {
        logoutFunc(accessToken);
      }
    },
  },
  computed: {
      isAuth() {
        return localStorage.isAuth === 'true';
      }
    }
}
</script>
<style scoped lang="scss">
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
@import "@/assets/styles/blocks/_main-nav.scss";
</style>
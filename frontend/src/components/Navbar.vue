<template>
  <nav class="main-nav">
    <div class="main-nav__wrapper">
      <ul class="main-nav__list">
        <router-link class="main-nav__logo" to="/"><img src="@/assets/logo.png" width="50" height="50" alt="Лого"/></router-link>
        <li class="main-nav__item">
          <router-link v-show="isAuth" class="main-nav__link main-nav__link--upload btn btn--white" to="/uploadmodel">
            <svg width="16px" height="16px" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 8L12 16" stroke="#323232" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M15 11L12.087 8.08704V8.08704C12.039 8.03897 11.961 8.03897 11.913 8.08704V8.08704L9 11" stroke="#323232" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M3 15L3 16L3 19C3 20.1046 3.89543 21 5 21L19 21C20.1046 21 21 20.1046 21 19L21 16L21 15" stroke="#323232" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span>Загрузка моделей</span>
          </router-link>
        </li>
        <li class="main-nav__item">
          <router-link v-if="!isAuth" class="main-nav__link btn" to="/authorization">Авторизация</router-link>
          <router-link v-if="isAuth" class="main-nav__link btn" to="/lk">
            <span class="main-nav__lk main-nav__lk--big">Личный кабинет</span>
            <span class="main-nav__lk main-nav__lk--small">ЛК</span>
          </router-link>
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
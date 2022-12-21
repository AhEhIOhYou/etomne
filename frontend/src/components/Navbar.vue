<template>
  <nav class="main-nav">
    <div class="main-nav__wrapper">
      <ul class="main-nav__list">
        <li class="main-nav__item">
          <a class="main-nav__link btn btn--white" href="/upload">Загрузка моделей</a>
          <!-- <router-link class="main-nav__link btn btn--white" to="/upload">Загрузка моделей</router-link> -->
        </li>
        <li class="main-nav__item">
          <!-- <a class="main-nav__link btn" href="/login">Авторизация</a> -->
          <router-link class="main-nav__link btn" to="/authorization">Авторизация</router-link>
          <button @click="unsetCookies" class="main-nav__link btn" type="button">Выйти</button>
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
  methods: {
    unsetCookies() {
      const accessToken = $cookies.get("access_token");
      $cookies.remove("access_token");
      $cookies.remove("refresh_token");
      console.log(accessToken);
      axios.get("/api/users/logout/", {
        headers: {
          "Authorization": `Bearer ${accessToken}`,
          // token: localStorage.getItem("access_token")
        }
      })
      .then(res => {
        console.log(res);
      });
    }
  }
}
</script>
<style scoped lang="scss">
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
@import "@/assets/styles/blocks/_main-nav.scss";
</style>
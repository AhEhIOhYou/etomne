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
          <router-link class="main-nav__link btn" to="/registration">Авторизация</router-link>
        </li>
      </ul>
    </div>
  </nav>
</template>
<script>
export default {  
  name: 'navbar'
}
</script>
<style scoped lang="scss">
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
@import "@/assets/styles/blocks/_main-nav.scss";
</style>
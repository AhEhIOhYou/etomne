<template>
  <div class="authorize">
    <div class="authorize__title">Пожалуйста, авторизуйтесь:</div>
    <div class="authorize__container">
      <div class="authorize__input-container input --grey">
        <label for="USER_LOGIN" class="authorize__label">Электронная почта:</label>
        <input class="authorize__input input__area" type="text" id="USER_LOGIN" autocomplete="off" name="USER_LOGIN" maxlength="255" placeholder="Введите почту">
      </div>
      <div class="authorize__input-container input --grey">
        <label for="USER_PASSWORD" class="authorize__label">Пароль:</label>
        <input class="authorize__input input__area" type="password" id="USER_PASSWORD" name="USER_PASSWORD" maxlength="255" autocomplete="off" placeholder="Введите пароль">
      </div>
    </div>
    <div class="authorize__button-container">
      <button type="submit" class="authorize__button btn">Войти</button>
      <router-link class="authorize__button btn btn--white" to="/registration">Зарегистрироваться</router-link>
    </div>
  </div>
</template>

<script>
export default {
  
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_authorize.scss";
@import "@/assets/styles/blocks/_alert.scss";
</style>
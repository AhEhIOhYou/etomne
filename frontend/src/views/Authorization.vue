<template>
  <div class="authorize">
    <div class="authorize__title">Пожалуйста, авторизуйтесь:</div>
    <div class="authorize__container">
      <div class="authorize__input-container input --grey">
        <label for="USER_LOGIN" class="authorize__label">Логин:</label>
        <input class="authorize__input input__area" type="text" id="USER_LOGIN" name="USER_LOGIN" maxlength="255" placeholder="Введите логин">
      </div>
      <div class="authorize__input-container input --grey">
        <label for="USER_PASSWORD" class="authorize__label">Пароль:</label>
        <input class="authorize__input input__area" type="password" id="USER_PASSWORD" name="USER_PASSWORD" maxlength="255" autocomplete="off" placeholder="Введите пароль">
      </div>
    </div>
    <div class="authorize__container">
      <div class="authorize__checkbox-container">
        <input class="visually-hidden authorize__checkbox" type="checkbox" id="USER_REMEMBER" name="USER_REMEMBER" value="Y">
        <label class="authorize__checkbox-label" for="USER_REMEMBER">Запомнить меня на этом компьютере</label>
      </div>
      <p>
        <a class="authorize__forget" href="/lk/?forgot_password=yes">Забыли свой пароль?</a>
      </p>
    </div>
    <div class="authorize__button-container">
      <button type="submit" class="authorize__button btn">Войти</button>
      <a class="authorize__button btn btn--white" href="/lk/?register=yes">Зарегистрироваться</a>
    </div>
  </div>
</template>

<script>
export default {
  
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_authorize.scss";
</style>
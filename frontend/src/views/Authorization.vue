<template>
  <div class="form">
    <form method="post" @submit.prevent>
      <p class="form__title">Пожалуйста, авторизуйтесь:</p>
      <p class="form__error alert alert--error">{{ $store.state.authorization.error }}</p>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="USER_EMAIL" class="form__label">Электронная почта:</label>
          <CustomInput :model-value="email" @update:model-value="setEmail" @input="inputListener" class="form__input" type="email" id="USER_EMAIL" name="USER_EMAIL" maxlength="255" placeholder="Введите почту" required/>
        </div>
      </div>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="USER_PASSWORD" class="form__label">Пароль:</label>
          <CustomInput :model-value="password" @update:model-value="setPassword" @input="inputListener" class="form__input" type="password" id="USER_PASSWORD" name="USER_PASSWORD" maxlength="255"  placeholder="Введите пароль" required/>
        </div>
      </div>
      <div class="form__button-container">
        <button type="submit" @click="handleSubmitAuthorization" class="form__button btn">Войти</button>
        <router-link class="form__button btn btn--white" to="/registration">Зарегистрироваться</router-link>
      </div>
    </form>
  </div>
</template>

<script>
import CustomInput from "@/components/UI/CustomInput";
import {mapState, mapMutations, mapActions} from 'vuex';

export default {
  components: {
    CustomInput
  },
  methods: {
    ...mapMutations({
      setEmail: 'authorization/setEmail',
      setPassword: 'authorization/setPassword',
      setError: 'authorization/setError'
    }),
    ...mapActions({
      handleSubmitAuthorization: 'authorization/handleSubmitAuthorization',
    }),
    inputListener() {
      const error = document.querySelector('.form__error');
      if (!error) {
        return;
      }
      if (error.textContent.length > 0) {
        error.classList.remove('alert--enable');
      }
    }
  },
  computed: {
    ...mapState({
      email: state => state.authorization.email,
      password: state => state.authorization.password,
      error: state => state.authorization.error
    }),
  }
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_form.scss";
@import "@/assets/styles/blocks/_alert.scss";
</style>
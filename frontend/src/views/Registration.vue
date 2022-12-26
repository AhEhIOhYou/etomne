<template>
  <div class="form">
    <form method="post" @submit.prevent>
      <p class="form__title">Пожалуйста, зарегистрируйтесь:</p>
      <p class="form__error alert alert--error">{{ $store.state.registration.error }}</p>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="USER_LOGIN" class="form__label">Имя пользователя:</label>
          <CustomInput :model-value="login" @update:model-value="setLogin" @input="inputListener" class="form__input" type="text" id="USER_LOGIN" name="USER_LOGIN" maxlength="255" placeholder="Введите имя пользователя" required/>
        </div>
      </div>
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
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="USER_CONFIRM" class="form__label">Подтверждение пароля:</label>
          <CustomInput :model-value="confirmPassword" @update:model-value="setConfirmPassword" @input="inputListener" class="form__input" type="password" id="USER_CONFIRM" name="USER_CONFIRM" maxlength="255"  placeholder="Введите пароль" required/>
        </div>
      </div>
      <button type="submit" @click="handleSubmitRegistration" class="form__button btn">Зарегистрироваться</button>
    </form>
  </div>
</template>

<script>
import {mapState, mapGetters, mapMutations, mapActions} from 'vuex';
import CustomInput from "@/components/UI/CustomInput";

export default {
  components: {
    CustomInput
  },
  methods: {
    ...mapMutations({
      setLogin: 'registration/setLogin',
      setEmail: 'registration/setEmail',
      setPassword: 'registration/setPassword',
      setConfirmPassword: 'registration/setConfirmPassword',
      setError: 'registration/setError'
    }),
    ...mapActions({
      handleSubmitRegistration: 'registration/handleSubmitRegistration',
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
      login: state => state.registration.login,
      email: state => state.registration.email,
      password: state => state.registration.password,
      confirmPassword: state => state.registration.confirmPassword,
      error: state => state.registration.error
    }),
  }
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_form.scss";
@import "@/assets/styles/blocks/_alert.scss";
</style>
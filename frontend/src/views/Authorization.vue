<template>
  <div class="authorize">
    <form method="post" @submit.prevent>
      <p class="authorize__title">Пожалуйста, авторизуйтесь:</p>
      <p class="authorize__error alert alert--error">{{ $store.state.authorization.error }}</p>
      <div class="authorize__container authorize__container--one">
        <div class="authorize__input-container input --grey">
          <label for="USER_EMAIL" class="authorize__label">Электронная почта:</label>
          <Input :model-value="email" @update:model-value="setEmail" @input="inputListener" class="authorize__input" type="email" id="USER_EMAIL" name="USER_EMAIL" maxlength="255" placeholder="Введите почту" required/>
        </div>
      </div>
      <div class="authorize__container authorize__container--one">
        <div class="authorize__input-container input --grey">
          <label for="USER_PASSWORD" class="authorize__label">Пароль:</label>
          <Input :model-value="password" @update:model-value="setPassword" @input="inputListener" class="authorize__input" type="password" id="USER_PASSWORD" name="USER_PASSWORD" maxlength="255"  placeholder="Введите пароль" required/>
        </div>
      </div>
      <div class="authorize__button-container">
        <button type="submit" @click="handleSubmitAuthorization" class="authorize__button btn">Войти</button>
        <router-link class="authorize__button btn btn--white" to="/registration">Зарегистрироваться</router-link>
      </div>
    </form>
  </div>
</template>

<script>
import {mapState, mapGetters, mapMutations, mapActions} from 'vuex';
import Input from "@/components/UI/Input";

export default {
  components: {
    Input
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
      const error = document.querySelector('.authorize__error');
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
@import "@/assets/styles/blocks/_authorize.scss";
@import "@/assets/styles/blocks/_alert.scss";
</style>
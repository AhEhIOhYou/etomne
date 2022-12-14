<template>
  <div class="authorize">
    <form method="post" @submit.prevent>
      <p class="authorize__title">Пожалуйста, зарегистрируйтесь:</p>
      <p class="authorize__error alert alert--error">{{ $store.state.registration.error }}</p>
      <div class="authorize__container authorize__container--one">
        <div class="authorize__input-container input --grey">
          <label for="USER_LOGIN" class="authorize__label">Имя пользователя:</label>
          <Input :model-value="login" @update:model-value="setLogin" class="authorize__input" type="text" id="USER_LOGIN" name="USER_LOGIN" maxlength="255" placeholder="Введите имя пользователя" required/>
        </div>
      </div>
      <div class="authorize__container authorize__container--one">
        <div class="authorize__input-container input --grey">
          <label for="USER_EMAIL" class="authorize__label">Электронная почта:</label>
          <Input :model-value="email" @update:model-value="setEmail" class="authorize__input" type="email" id="USER_EMAIL" name="USER_EMAIL" maxlength="255" placeholder="Введите почту" required/>
        </div>
      </div>
      <div class="authorize__container authorize__container--one">
        <div class="authorize__input-container input --grey">
          <label for="USER_PASSWORD" class="authorize__label">Пароль:</label>
          <Input :model-value="password" @update:model-value="setPassword" class="authorize__input" type="password" id="USER_PASSWORD" name="USER_PASSWORD" maxlength="255"  placeholder="Введите пароль" required/>
        </div>
      </div>
      <div class="authorize__container authorize__container--one">
        <div class="authorize__input-container input --grey">
          <label for="USER_CONFIRM" class="authorize__label">Подтверждение пароля:</label>
          <Input :model-value="confirmPassword" @update:model-value="setConfirmPassword" class="authorize__input" type="password" id="USER_CONFIRM" name="USER_CONFIRM" maxlength="255"  placeholder="Введите пароль" required/>
        </div>
      </div>
      <button type="submit" @click="handleSubmit" class="authorize__button btn">Зарегистрироваться</button>
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
      setLogin: 'registration/setLogin',
      setEmail: 'registration/setEmail',
      setPassword: 'registration/setPassword',
      setConfirmPassword: 'registration/setConfirmPassword',
      setError: 'registration/setError'
    }),
    ...mapActions({
      handleSubmit: 'registration/handleSubmit',
    }),
  },
  computed: {
    ...mapState({
      login: state => state.registration.login,
      email: state => state.registration.email,
      password: state => state.registration.password,
      confirmPassword: state => state.registration.confirmPassword,
      error: state => state.registration.error
    }),
    ...mapGetters({
      inputText: 'registration/inputText'
    })
  }
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_authorize.scss";
@import "@/assets/styles/blocks/_alert.scss";
</style>
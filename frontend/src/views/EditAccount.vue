<template>
  <div class="form">
    <form method="post" @submit.prevent>
      <p class="form__title">Введите изменяемые данные:</p>
      <p class="form__error alert alert--error">{{ $store.state.edit.error }}</p>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="USER_LOGIN" class="form__label">Имя пользователя:</label>
          <CustomInput :model-value="login" @update:model-value="setLogin" @input="inputListener" ref="username" class="form__input" type="text" id="USER_LOGIN" name="USER_LOGIN" maxlength="255" placeholder="Введите новое имя пользователя" required/>
        </div>
      </div>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="USER_EMAIL" class="form__label">Электронная почта:</label>
          <CustomInput :model-value="email" @update:model-value="setEmail" @input="inputListener" ref="useremail" class="form__input" type="email" id="USER_EMAIL" name="USER_EMAIL" maxlength="255" placeholder="Введите новую почту" required/>
        </div>
      </div>
      <button type="submit" @click="handleSubmitEdit" class="form__button btn">Обновить данные</button>
      <router-link class="form__button btn btn--white" to="/lk">Отмена</router-link>
    </form>
  </div>
</template>

<script>
import {mapState, mapMutations, mapActions} from 'vuex';
import axios from 'axios';
import CustomInput from "@/components/UI/CustomInput";

export default {
  components: {
    CustomInput
  },
  methods: {
    ...mapMutations({
      setLogin: 'edit/setLogin',
      setEmail: 'edit/setEmail',
      setError: 'edit/setError'
    }),
    ...mapActions({
      handleSubmitEdit: 'edit/handleSubmitEdit',
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
      login: state => state.edit.login,
      email: state => state.edit.email,
      error: state => state.edit.error
    }),
  },
  mounted () {
    const accessToken = $cookies.get("access_token");
    const refreshToken = $cookies.get("refresh_token");
    const userId = localStorage.getItem("id");

    const showUserInfo = (id) => {
      axios.get(`api/users/${id}`
      ).then(response => {
          this.userInfo = response.data;
          this.$store.commit('edit/setLogin', this.userInfo.name);
          this.$store.commit('edit/setEmail', this.userInfo.email);
        })
        .catch(error => {
          console.log(error);
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
          showUserInfo(userId);
        })
        .catch(error => {
          console.log(error);
        });
    } else {
      showUserInfo(userId);
    }
  },
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_form.scss";
@import "@/assets/styles/blocks/_alert.scss";
</style>
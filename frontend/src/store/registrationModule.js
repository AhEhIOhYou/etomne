import axios from "axios";
import VueCookies from 'vue-cookies';

export const registrationModule = {
    state: () => ({
        login: '',
        email: '',
        password: '',
        confirmPassword: '',
        error: '',
    }),
    mutations: {
      setLogin(state, login) {
        state.login = login;
      },
      setEmail(state, email) {
        state.email = email; 
      },
      setPassword(state, password) {
        state.password = password;
      },
      setConfirmPassword(state, confirmPassword) {
        state.confirmPassword = confirmPassword;
      },
      setError(state, error) {
        state.error = error;
      },
    },
    getters: {
    },
    actions: {
      handleSubmitRegistration({state, commit}) {
        const error = document.querySelector('.alert--error');
        if(!(state.password === state.confirmPassword)) {
          commit('setError', 'Пароли не совпадают');
          error.classList.add('alert--enable');
        } else {
          commit('setError', '');
          error.classList.remove('alert--enable');
          axios.post('/api/users', {
            name: state.login,
            email: state.email,
            password: state.password
          })
          .then(response => {
            window.location.href = '/';
            console.log(response);
          })
          .catch(error => {
            console.log(error);
          });
        }
      },
    },
    namespaced: true
}
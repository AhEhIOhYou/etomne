import axios from "axios";
import VueCookies from 'vue-cookies';

export const authorizationModule = {
    state: () => ({
        email: '',
        password: '',
        error: '',
        id: '',
        name: ''
    }),
    mutations: {
      setEmail(state, email) { 
        state.email = email;
      },
      setPassword(state, password) {
        state.password = password;
      },
      setError(state, error) {
        state.error = error;
      },
      setId(state, id) {
        state.id = id;
      }, 
      setName(state, name) {
        state.name = name;
      },
    },
    getters: { 
    },
    actions: {
      handleSubmitAuthorization({state, commit}) {
        const errorElem = document.querySelector('.alert--error');
        axios.post('/api/users/login', {
          email: state.email,
          password: state.password
        })
        .then(response => {
          console.log(response);
          $cookies.set('access_token', response.data.tokens.access_token, '15min', '/');
          $cookies.set('refresh_token', response.data.tokens.refresh_token, '7d', '/');
          commit('setId', `${response.data.public_data.id}`);
          commit('setName', `${response.data.public_data.name}`);
          localStorage.setItem('name', response.data.public_data.name);
          localStorage.setItem('id', response.data.public_data.id);
          localStorage.setItem('isAuth', true);
          window.location.href = '/';
        })
        .catch(error => {
          errorElem.classList.add('alert--enable');
          commit('setError', 'Неверные данные');
        });
      },
    },
    namespaced: true
}
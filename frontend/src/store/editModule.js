import axios from "axios";
import VueCookies from 'vue-cookies';

export const editModule = {
    state: () => ({
        login: '',
        email: '',
        error: '',
    }),
    mutations: {
      setLogin(state, login) {
        state.login = login;
      },
      setEmail(state, email) {
        state.email = email; 
      },
      setError(state, error) {
        state.error = error;
      },
    },
    getters: {
    },
    actions: {
      handleSubmitEdit({state, commit}) {
        const userId = localStorage.getItem("id");
        const accessToken = $cookies.get("access_token");
        const error = document.querySelector('.alert--error');
        if(state.login.length === 0) {
          commit('setError', 'Имя пользователя не может быть пустым');
          error.classList.add('alert--enable');
        } else if (state.email.length === 0) {
          commit('setError', 'Почта не может быть пустой');
          error.classList.add('alert--enable'); 
        } else {
          commit('setError', '');
          error.classList.remove('alert--enable');
          const data = {
            name: state.login,
            email: state.email,
          }
          axios.post(`/api/users/update/${userId}`,
          data,
          { 
            headers: { 
              "Authorization": `Bearer ${accessToken}`
            } 
          })
          .then(response => {
            window.location.href = '/lk';
          })
          .catch(error => {
            console.log(error);
          });
        }
      },
    },
    namespaced: true
}
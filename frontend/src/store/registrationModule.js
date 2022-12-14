export const registrationModule = {
    state: () => ({
        login: '',
        email: '',
        password: '',
        confirmPassword: '',
        error: 'АБОБА',
    }),
    mutations: {
      setLogin(state, login) {
        state.login = login;
      },
      setEmail(state, email) {
        state.email = email
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
      handleSubmit({state, commit}) {
        // Если пароли не совпадают(comparePasswords), то поменять state.error на 'Пароли не совпадают' и вставить в значение <p class="authorize__error alert alert--error"></p>
        const error = document.querySelector('.alert--error');
        console.log(state.password);
        console.log(state.confirmPassword);
        if(!(state.password === state.confirmPassword)) {
          commit('setError', 'Пароли не совпадают');
          error.classList.add('alert--enable');
        } else {
          commit('setError', '');
          error.classList.remove('alert--enable');
        }
      }
    },
    namespaced: true
}
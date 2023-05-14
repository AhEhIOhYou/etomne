import axios from "axios";
import VueCookies from 'vue-cookies';

export const uploadModule = {
    state: () => ({
        name: '',
        description: ''
    }),
    mutations: {
      setName(state, name) {
        state.name = name;
      },
      setDescription(state, description) {
        state.description = description;
      },
    },
    getters: { 
    },
    actions: {
      handleSubmitUpload({state, commit}) {
        axios.post('/api/model', {
          name: state.email,
          description: state.password
        })
        .then(response => {
        })
        .catch(error => {
          console.loog(error);
        });
      },
    },
    namespaced: true
}
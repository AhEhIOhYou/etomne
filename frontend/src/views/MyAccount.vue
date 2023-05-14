<template>
  <div class="account">
    <dl class="account__info-list">
      <dt class="account__info-term">Имя пользователя:</dt>
      <dd class="account__info-description">{{ this.userInfo.name}}</dd>
      <dt class="account__info-term">Почта:</dt>
      <dd class="account__info-description">{{ this.userInfo.email }}</dd>
      <dt class="account__info-term">Является ли админом:</dt>
      <dd v-if="this.userInfo.is_admin" class="account__info-description">Является</dd>
      <dd v-else class="account__info-description">Не является</dd>
    </dl>
    <router-link class="account__btn btn" to="/edit">Изменить данные</router-link>
  </div>
  <section class="models">
    <models-list
      :models="models"
      @remove="removeModel"
      v-if="!isModelsLoading"
    />
    <div v-else>Идет загрузка...</div>
    <div v-intersection="loadMoreModels" class="observer"></div>
</section>
</template>

<script>
import {mapState, mapActions, mapMutations} from 'vuex';
import axios from "axios";

export default {
  components: {
  },
  data() {
    return {
      userInfo: null,
    } 
  },
  methods: {
    ...mapMutations({
      setPage: 'models/setPage',
      setUserId: 'models/setUserId'
    }),
    ...mapActions({
      loadMoreModels: 'models/loadMoreModels',
      fetchModels: 'models/fetchModels',
      setPagesToOne: 'models/setPagesToOne',
    }),
    removeModel(model){
      const accessToken = $cookies.get("access_token");
      const modelId = model.model.id;
      axios.delete(`/api/model/${modelId}`,
      { 
        data: { 
          id: modelId
        }, 
        headers: { 
          "Authorization": `Bearer ${accessToken}`
        } 
      })
      .then(response => {
        const model = document.querySelector(`[data-model-id="${modelId}"`);
        model.remove();
        const models = document.querySelectorAll('.model');
        if(models.length <= 1) {
          this.loadMoreModels();
        }
      })
      .catch(error => {
        console.log(error);
      });
    }
  },
  props: {
  },
  computed: {
    ...mapState({
      models: state => state.models.models,
      isModelsLoading: state => state.models.isModelsLoading,
      page: state => state.models.page,
      limit: state => state.models.limit,
      totalPages: state => state.models.totalPages,
    }),
  },
  mounted() {
    this.setPagesToOne();
    this.fetchModels();
    let modelViewerScript = document.createElement('script');
    modelViewerScript.setAttribute('src', 'https://unpkg.com/@google/model-viewer/dist/model-viewer.min.js');
    modelViewerScript.setAttribute('type', 'module');
    document.head.appendChild(modelViewerScript);
  },
  created() {
    const accessToken = $cookies.get("access_token");
    const refreshToken = $cookies.get("refresh_token");
    const userId = localStorage.getItem("id");
    this.$store.commit('models/setUserId', userId);

    const showUserInfo = (id) => {
      axios.get(`api/users/${id}`
      ).then(response => {
          this.userInfo = response.data;
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
@import "@/assets/styles/blocks/_account.scss";

.observer {
  height: 30px;
}
</style>
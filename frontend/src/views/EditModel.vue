<template>
  <div v-if="!isModelLoading" class="form">
    <div class="form__overlay">
      <div class="form__spinner" role="status">
        <span class="form__loading">Загрузка...</span>
      </div>
    </div>
    <form class="form__form" method="post" @submit.prevent>
      <p class="form__title">Пожалуйста, введите данные для редактирования модели:</p>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="MODEL_NAME" class="form__label">Новое название модели:</label>
          <CustomInput :model-value="name" @update:model-value="setName" class="form__input" type="text" id="MODEL_NAME" name="MODEL_NAME" maxlength="255" :value="this.model.model.title" placeholder="Новое название модели" required/>
        </div>
      </div>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="MODEL_DESCRIPTION" class="form__label">Новое описание модели:</label>
          <CustomTextarea :model-value="description" @update:model-value="setDescription" class="form__input" id="MODEL_DESCRIPTION" name="MODEL_DESCRIPTION" :value="this.model.model.description" placeholder="Новое описание модели"/>
        </div>
      </div>
      <EditFiles :model="model" @onChange='onSaved'/>
      <button type="submit" @click="submitFilesHandler" class="form__button btn">Обновить модель</button>
    </form>
    <div class="modal">
      <p class="modal__text">Модель успешно обновлена</p>
    </div>
  </div>
  <div v-else>Идет загрузка...</div> 
</template>

<script>
import {mapState, mapGetters, mapMutations, mapActions} from 'vuex';
import axios from "axios";
import EditFiles from "@/components/EditFiles.vue";
import CustomInput from "@/components/UI/CustomInput";
import CustomTextarea from "@/components/UI/CustomTextarea";

export default {
  components: {
    CustomInput,
    CustomTextarea,
    EditFiles
  },
  data(){
    return {
      model: null,
      isModelLoading: true,
      attachments: '',
      files_id: []
    }
  },
  methods: {
    ...mapMutations({
      setName: 'upload/setName',
      setDescription: 'upload/setDescription',
    }),
    onSaved (data) {
      this.files_id = data;
    }, 
    handleFilesUpload(){
      this.attachments = this.$refs.attachments.files;
    },
    submitFilesHandler(evt){
      evt.preventDefault();
      const accessToken = $cookies.get("access_token");
      const refreshToken = $cookies.get("refresh_token");
      const formOverlay = document.querySelector('.form__overlay');
      const ids = JSON.parse(JSON.stringify(this.files_id));
      const modelId = this.model.model.id;

      const title = this.name;
      const description = this.description;
      const data = {
        description: description,
        title: title,
      }

      const submitFilesFunc = (model_id, access) => {
      axios.post(`/api/model/update/${model_id}`,
        data,
        {
          headers: {
            'Authorization': `Bearer ${access}`
          }
        }
      ).then(response => {
          formOverlay.classList.remove('form__overlay--active');
          const form = document.querySelector('.form__form');
          form.reset();
          const modal = document.querySelector('.modal');
          modal.classList.add('modal--active');
          setTimeout(() => {
            modal.classList.remove('modal--active');
          }, 5000);
        })
        .catch(error => {
          formOverlay.classList.remove('form__overlay--active');
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
            formOverlay.classList.add('form__overlay--active');
            submitFilesFunc(modelId, response.data.tokens.access_token);
            window.location.href = '/';
          })
          .catch(error => {
            console.log(error);
          });
      } else {
        formOverlay.classList.add('form__overlay--active');
        submitFilesFunc(modelId, accessToken);
      }
    },
    async fetchModel(id) {
      axios
      .get(`/api/model/${id}`)
      .then(response => {
        this.model = response.data;
      })
      .catch(error => {
        console.log(error);
      })
      .finally(() => (this.isModelLoading = false));
    },
  },
  computed: {
    ...mapState({
      name: state => state.upload.name,
      description: state => state.upload.description,
    }),
  },
  mounted() {
    this.fetchModel(this.$route.params.id);
  }
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_form.scss";
@import "@/assets/styles/blocks/_alert.scss";
@import "@/assets/styles/blocks/_modal.scss";
</style>
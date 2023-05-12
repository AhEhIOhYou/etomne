<template>
  <div class="form">
    <div class="form__overlay">
      <div class="form__spinner" role="status">
        <span class="form__loading">Загрузка...</span>
      </div>
    </div>
    <form class="form__form" method="post" @submit.prevent>
      <p class="form__title">Пожалуйста, введите данные для загрузки модели:</p>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="MODEL_NAME" class="form__label">Название модели:</label>
          <CustomInput :model-value="name" @update:model-value="setName" class="form__input" type="text" id="MODEL_NAME" name="MODEL_NAME" maxlength="255" placeholder="Введите название модели" required/>
        </div>
      </div>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="MODEL_DESCRIPTION" class="form__label">Описание модели:</label>
          <CustomTextarea :model-value="description" @update:model-value="setDescription" class="form__input" id="MODEL_DESCRIPTION" name="MODEL_DESCRIPTION" placeholder="Введите описание модели"/>
        </div>
      </div>
      <UploadFiles @onChange='onSaved'/>
      <button type="submit" @click="submitFilesHandler" class="form__button btn">Загрузить модель</button>
    </form>
    <div class="modal">
      <p class="modal__text">Модель успешно загружена</p>
    </div>
  </div>
</template>

<script>
import {mapState, mapGetters, mapMutations, mapActions} from 'vuex';
import axios from "axios";
import UploadFiles from "@/components/UploadFiles.vue";
import CustomInput from "@/components/UI/CustomInput";
import CustomTextarea from "@/components/UI/CustomTextarea";

export default {
  components: {
    CustomInput,
    CustomTextarea,
    UploadFiles
  },
  data(){
    return {
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

      const title = this.name;
      const description = this.description;
      const data = {
        description: description,
        title: title,
        files_id: ids.files_id
      }

      const submitFilesFunc = (access) => {
      axios.post('/api/model',
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
          window.location.href = '/authorization';
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
            submitFilesFunc(response.data.tokens.access_token);
            console.log(response);
          })
          .catch(error => {
            console.log(error);
          });
      } else {
        formOverlay.classList.add('form__overlay--active');
        submitFilesFunc(accessToken);
      }
    },
  },
  computed: {
    ...mapState({
      name: state => state.upload.name,
      description: state => state.upload.description,
    }),
  }
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_form.scss";
@import "@/assets/styles/blocks/_alert.scss";
@import "@/assets/styles/blocks/_modal.scss";
</style>
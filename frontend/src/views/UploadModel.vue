<template>
  <div class="form">
    <div class="form__overlay">
      <div class="form__spinner" role="status">
        <span class="form__loading">Loading...</span>
      </div>
    </div>
    <form class="form__form" method="post" @submit.prevent>
      <p class="form__title">Пожалуйста, введите данные для загрузки модели:</p>
      <!-- <p class="form__error alert alert--error">{{ $store.state.authorization.error }}</p> -->
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
      <div class="form__container form__container--one">
        <div class="form__upload-file">
          <p class="form__upload-text">Загрузить файлы:</p>
          <input class="visually-hidden form__upload-input" ref="attachments" v-on:change="handleFilesUpload" type="file" name="MODEL_FILES" id="MODEL_FILES" accept="image/*, video/*, .glb, .gltf" multiple required>
          <label class="form__upload-label" for="MODEL_FILES">
            <span class="btn form__upload-button">
              <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" viewBox="0 0 25 25" fill="none">
                <path d="M6.25013 18.7501H18.7501V20.8334H6.25013V18.7501ZM12.5001 3.73547L5.51367 10.7219L6.98659 12.1948L11.4585 7.72297V16.6667H13.5418V7.72297L18.0137 12.1948L19.4866 10.7219L12.5001 3.73547Z"/>
              </svg>
              <span>Выбрать файлы</span>
            </span>
          </label>
        </div>
      </div>
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
import CustomInput from "@/components/UI/CustomInput";
import CustomTextarea from "@/components/UI/CustomTextarea";

export default {
  components: {
    CustomInput,
    CustomTextarea
  },
  data(){
    return {
      attachments: '',
    }
  },
  methods: {
    ...mapMutations({
      setName: 'upload/setName',
      setDescription: 'upload/setDescription',
    }),
    ...mapActions({
      handleSubmitUpload: 'upload/handleSubmitUpload',
    }),
    handleFilesUpload(){
      this.attachments = this.$refs.attachments.files;
    },
    submitFilesHandler(){
      const accessToken = $cookies.get("access_token");
      const refreshToken = $cookies.get("refresh_token");
      const formOverlay = document.querySelector('.form__overlay');
      let formData = new FormData();
      formData.append('title', this.name);
      formData.append('description', this.description);
      for( var i = 0; i < this.attachments.length; i++ ){
          let attachment = this.attachments[i];
          formData.append('attachments', attachment);
        }

      const submitFilesFunc = (access) => {
      axios.post( '/api/model',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
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
          console.log(response);
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
            $cookies.set('access_token', response.data.access_token, '15min', '/');
            $cookies.set('refresh_token', response.data.refresh_token, '7d', '/');
            localStorage.setItem('isAuth', true);
            formOverlay.classList.add('form__overlay--active');
            submitFilesFunc(response.data.access_token);
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
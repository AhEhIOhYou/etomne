<template>
  <div class="form">
    <form method="post" @submit.prevent>
      <p class="form__title">Пожалуйста, введите данные для загрузки модели:</p>
      <!-- <p class="form__error alert alert--error">{{ $store.state.authorization.error }}</p> -->
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="MODEL_NAME" class="form__label">Название модели:</label>
          <CustomInput :model-value="email" @update:model-value="setEmail" @input="inputListener" class="form__input" type="email" id="MODEL_NAME" name="MODEL_NAME" maxlength="255" placeholder="Введите название модели" required/>
        </div>
      </div>
      <div class="form__container form__container--one">
        <div class="form__input-container input --grey">
          <label for="MODEL_DESCRIPTION" class="form__label">Описание модели:</label>
          <CustomTextarea :model-value="password" @update:model-value="setPassword" @input="inputListener" class="form__input" id="MODEL_DESCRIPTION" name="MODEL_DESCRIPTION" placeholder="Введите описание модели" required/>
        </div>
      </div>
      <div class="form__container form__container--one">
        <div class="form__upload-file">
          <p class="form__upload-text">Загрузить файлы:</p>
          <input class="visually-hidden form__upload-input" type="file" name="MODEL_FILES" id="MODEL_FILES" accept=".glb, .png, .jpeg, .jpg" multiple required>
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
      <button type="submit" @click="" class="form__button btn">Загрузить модель</button>
    </form>
  </div>
</template>

<script>
import {mapState, mapGetters, mapMutations, mapActions} from 'vuex';
import CustomInput from "@/components/UI/CustomInput";
import CustomTextarea from "@/components/UI/CustomTextarea";

export default {
  components: {
    CustomInput,
    CustomTextarea
  },
  methods: {
    ...mapMutations({
      setEmail: 'authorization/setEmail',
      setPassword: 'authorization/setPassword',
      setError: 'authorization/setError'
    }),
    ...mapActions({
      handleSubmitAuthorization: 'authorization/handleSubmitAuthorization',
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
      email: state => state.authorization.email,
      password: state => state.authorization.password,
      error: state => state.authorization.error
    }),
  }
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/blocks/_form.scss";
@import "@/assets/styles/blocks/_alert.scss";
</style>
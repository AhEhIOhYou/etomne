<template>
  <div>
    <div v-if="!isModelLoading" class="model">
      <div class="model__content">
        <h2 class="model__title">{{ model.model.title }}</h2>
        <model-viewer class="model__model" :src="'https://modelshowtime.serdcebolit.ru/' + model.files.glb[0].url" camera-controls="" ar-status="not-presenting"></model-viewer>
        <span class="model__author">Created by {{ model.author.name }}</span>
        <span class="model__data">{{ model.model.created_at }}</span>
      </div>
      <div class="model__panel">
        <div class="model__info">
          <h3 class="model__sub-title">Описание</h3>
          <div class="model__info-container">
            <p class="model__description">{{ model.model.description }}</p>
            <ul class="model__actions">
              <li class="model__action">
                <button class="model__action-btn btn">Редактировать</button>
              </li>
              <li class="model__action">
                <button class="model__action-btn btn">Удалить</button>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
    <div v-else>Идет загрузка...</div> 

    <!-- <div v-if="!isModelLoading" class="model__content">
      <h2 class="model__title">{{ model.model.title }}</h2>
      <img 
      v-if="model.files.length > 0"
      v-for="file in model.files"
      :src="'https://modelshowtime.serdcebolit.ru/' + file.url"
      >
      <img v-else src="https://placehold.co/600x400">
      <span class="model__author">{{ model.model.name }}</span>
      <span class="model__data">{{ model.model.created_at}}</span>
      <span class="model__description">{{ model.model.description }}</span>
    </div>
    <div v-else>Идет загрузка...</div>  -->
    </div>
</template>

<script>
import axios from "axios";
export default {
  components: {
  },
  data() {
    return {
      model: null,
      isModelLoading: true,
    }
  },
  methods: {
    async fetchModel(id) {
          axios
          .get(`https://modelshowtime.serdcebolit.ru/api/model/${id}`)
          .then(response => {
            this.model = response.data
          })
          .catch(error => {
            console.log(error);
          })
          .finally(() => (this.isModelLoading = false));
        },
  },
  mounted() {
    this.fetchModel(this.$route.params.id);
    let modelViewerScript = document.createElement('script')
    modelViewerScript.setAttribute('src', 'https://unpkg.com/@google/model-viewer/dist/model-viewer.min.js')
    modelViewerScript.setAttribute('type', 'module')
    document.head.appendChild(modelViewerScript)
  },
  computed: {
  },
}
</script>

<style lang="scss" scoped>
</style>
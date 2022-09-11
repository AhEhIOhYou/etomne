<template>
  <div>
    <h1>Это страница поста с ID = {{ $route.params.id }}</h1>
    <!-- <span> {{ modelFiles[0] }}</span> -->
    <!-- <img src="https://modelshowtime.serdcebolit.ru/" {{ model.files[0].url }}> -->
    <div v-if="!isModelLoading" class="model__content">
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
    <div v-else>Идет загрузка...</div>
  </div>
</template>

<script>
import {mapState, mapGetters, mapActions, mapMutations} from 'vuex';
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
    // ...mapActions({
    //   fetchModel: 'model/fetchModel',
    // }),
    async fetchModel(id) {
          //   try {
          //       commit('setLoading', true);
          //       const response = await axios.get(`https://modelshowtime.serdcebolit.ru/api/model/${id}`);
          //       commit('setModel', response.data);
          //   } catch (e) {
          //       console.log(e)
          //   } finally {
          //     commit('setLoading', false);
          // }
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
    
  },
  computed: {
    // ...mapState({
    //   model: state => state.model.model,
    //   isModelLoading: state => state.model.isModelLoading,
    // }),
  },
}
</script>

<style scoped>
</style>

<!-- <model-viewer class="model__model" src={{ model.FilePath }} camera-controls></model-viewer> -->
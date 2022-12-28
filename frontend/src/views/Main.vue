<template>
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
import {mapState, mapGetters, mapActions, mapMutations} from 'vuex';
import axios from "axios";

export default {
  components: {
  },
  data() {
    return {
    }
  },
  methods: {
    ...mapMutations({
      setPage: 'models/setPage',
    }),
    ...mapActions({
      loadMoreModels: 'models/loadMoreModels',
      fetchModels: 'models/fetchModels'
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
        console.log(response);
      })
      .catch(error => {
        console.log(error);
      });
    }
  },
  mounted() {
    this.fetchModels();
    let modelViewerScript = document.createElement('script');
    modelViewerScript.setAttribute('src', 'https://unpkg.com/@google/model-viewer/dist/model-viewer.min.js');
    modelViewerScript.setAttribute('type', 'module');
    document.head.appendChild(modelViewerScript);
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
}
</script>

<style lang="scss" scoped>

</style>
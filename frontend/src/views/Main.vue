<template>
  <section class="models">
      <models-list
        :models="models"
        v-if="!isModelsLoading"
      />
      <div v-else>Идет загрузка...</div>
      <div v-intersection="loadMoreModels" class="observer"></div>
  </section>
</template>

<script>
import {mapState, mapGetters, mapActions, mapMutations} from 'vuex'
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
  },
  mounted() {
    this.fetchModels();
    let modelViewerScript = document.createElement('script')
    modelViewerScript.setAttribute('src', 'https://unpkg.com/@google/model-viewer/dist/model-viewer.min.js')
    modelViewerScript.setAttribute('type', 'module')
    document.head.appendChild(modelViewerScript)
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
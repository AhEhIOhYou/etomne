<template>
  <div class="model">
    <div class="model__content">
      <h2 class="model__title">{{ model.model.title }}</h2>
      <!-- <model-viewer class="model__model" src={{ model.FilePath }} camera-controls></model-viewer> -->
      <img 
      v-if="model.files.length > 0"
      v-for="file in model.files"
      :src="'https://modelshowtime.serdcebolit.ru/' + file.url"
      >
      <img v-else src="https://placehold.co/600x400">
      <button @click="$router.push(`/${model.model.id}`)"></button>
      <span class="model__author">{{ model.author.name }}</span>
      <span class="model__data">{{ model.model.created_at}}</span>
      <span class="model__description">{{ model.model.description }}</span>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    model: {
      type: Object,
      required: true,
    }
  },
  name: 'models-item',
}
</script>

<style lang="scss" scoped>
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";

.model {
  &__title {
    margin: 0 0 8px;
    font-size: 18px;
    font-weight: 600;
    line-height: 22px;
    color: $color-rustic-red;
  }

  &__model {
    background-color: $color-default-white;
    width: 100%;
    height: 512px;
    margin-bottom: 10px;
  }

  &__author,
  &__data,
  &__comment-author {
    display: block;
    color: $color-regent-grey;
    font-size: 12px;
    font-weight: 400;
    line-height: 16px;
  }
}
</style>
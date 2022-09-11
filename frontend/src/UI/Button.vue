<template>
  <button class="btn">
    <slot></slot>
  </button>
</template>

<script>
export default {
  name: 'default-button'
}
</script>

<style lang="scss" scoped>
  @import "@/assets/styles/_variables.scss";
  @import "@/assets/styles/_mixins.scss";
  @import "@/assets/styles/blocks/_btn.scss";
</style>
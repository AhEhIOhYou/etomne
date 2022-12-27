<template>
  <textarea :value="modelValue" @input="updateInput" class="input input--textarea" autocomplete="off"></textarea>
</template>

<script>
export default {
  name: 'CustomTextarea',
  props: {
    modelValue: [String, Number]
  },
  methods: {
    updateInput(event) {
      this.$emit('update:modelValue', event.target.value)
    }
  }
}
</script>

<style scoped lang="scss">
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
@import "@/assets/styles/blocks/_input.scss";
</style>

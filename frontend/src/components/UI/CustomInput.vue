<template>
  <input :value="modelValue" @input="updateInput" class="input" autocomplete="off" type="text">
</template>

<script>
export default {
  name: 'CustomInput',
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
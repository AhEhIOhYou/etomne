<template>
  <input :value="modelValue" @input="updateInput" class="input__area" autocomplete="off" type="text">
</template>

<script>
export default {
  name: 'Input',
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

<style scoped>


</style>

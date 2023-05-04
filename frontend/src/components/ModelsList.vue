<template>
  <div class="models-list">
    <transition-group name="models-list">
      <models-item
        v-for="model in models"
        :model="model"
        :key="model.model.id"
        @remove="$emit('remove', model)"
      />
    </transition-group>
  </div>
</template>

<script>
export default {
  props: {
    models: {
      type: Array,
      required: true
    }
  },
  name: 'models-list'
}
</script>

<style lang="scss" scoped> 
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";

.models-list-enter-active,
.models-list-leave-active {
  transition: all 0.4s ease;
}
.models-list-enter-from,
.models-list-leave-to {
  opacity: 0;
  transform: translateX(130px);
}
.models-list-move {
  transition: transform 0.4s ease;
}
</style>
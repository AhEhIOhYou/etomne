<template>
  <header class="header">
    <navbar></navbar>
  </header>
  <div class="container">
    <router-view/>
  </div>
</template>

<style lang="scss">
@import "@/assets/styles/vendor/_normalize.scss";
@import "@/assets/styles/_variables.scss";
@import "@/assets/styles/_mixins.scss";
@import "@/assets/styles/global/_fonts.scss";
@import "@/assets/styles/global/_reboot.scss";
@import "@/assets/styles/global/_utils.scss";
@import "@/assets/styles/global/_container.scss";
@import "@/assets/styles/blocks/_header.scss";
@import "@/assets/styles/blocks/_btn.scss";
@import "@/assets/styles/blocks/_model.scss";
</style>

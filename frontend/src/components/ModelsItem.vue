<template>
  <!-- <div class="model">
    <div class="model__content">
      <h2 class="model__title">{{ model.model.title }}</h2>
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
  </div> -->

  <div class="model">
    <div class="model__content">
      <h2 class="model__title">{{ model.model.title }}</h2>
      <span>{{model.model.id}}</span>
      <div class="model__swipers">
        <swiper
          :slides-per-view="1"
          :space-between="20"
          :allowTouchMove="false"
          :modules="[Thumbs]"
          :thumbs="{ swiper: thumbsSwiper }"
          class="model__default-swiper"
        >
          <swiper-slide v-for="model in model.files.glb">
            <model-viewer class="model__model" :src="'https://modelshowtime.serdcebolit.ru/' + model.url" powerPreference="low-power" camera-controls=""></model-viewer>
          </swiper-slide>
          <swiper-slide v-for="img in model.files.img">
            <img 
              :src="'https://modelshowtime.serdcebolit.ru/' + img.url"
            >
          </swiper-slide>
        </swiper>
        <swiper v-if="model.files.glb.length > 1 || model.files.img.length > 0" class="model__thumbs-swiper"
          @swiper="setThumbsSwiper"
          :slides-per-view="4"
          :space-between="10"
          :freeMode="true"
          :watchSlidesProgress="true"
          :grabCursor="true"
          :modules="[Thumbs, FreeMode, Navigation]"
          :mousewheel="true"
          :navigation="true"
        >
          <swiper-slide v-for="model in model.files.glb">
            <model-viewer class="model__model" :src="'https://modelshowtime.serdcebolit.ru/' + model.url"></model-viewer>
          </swiper-slide>
          <swiper-slide v-for="img in model.files.img">
            <img 
              :src="'https://modelshowtime.serdcebolit.ru/' + img.url"
            >
          </swiper-slide>
        </swiper>
      </div>
      <!-- <div v-if="model.files.glb.length > 0">
        <model-viewer class="model__model" :src="'https://modelshowtime.serdcebolit.ru/' + model.files.glb[0].url" camera-controls="" ar-status="not-presenting"></model-viewer>
      </div>
      <span v-else>Модельки нет</span>
      <img 
      v-for="img in model.files.img"
      :src="'https://modelshowtime.serdcebolit.ru/' + img.url"
      > -->
      <span class="model__author">Created by {{ model.author.name }}</span>
      <!-- https://modelshowtime.serdcebolit.ru/upload/f391b5a497b187f970c3f9bf9ad999a9.glb" -->
      <span class="model__data">{{ model.model.created_at }}</span>
      <button class="btn" @click="$router.push(`/${model.model.id}`)">Перейти к подробному описанию модели</button>
    </div>
    <div class="model__panel">
      <div class="model__info">
        <h3 class="model__sub-title">Описание</h3>
        <div class="model__info-container">
          <p class="model__description">{{ model.model.description }}</p>
          <ul class="model__actions">
            <li class="model__action">
              <button class="model__action-btn btn">Редактировать</button>
            </li>
            <li class="model__action">
              <button class="model__action-btn btn">Удалить</button>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import { Swiper, SwiperSlide } from 'swiper/vue';
import 'swiper/css';
import "swiper/css/free-mode"
import "swiper/css/thumbs"
import "swiper/css/navigation"
import {FreeMode, Thumbs, Navigation} from 'swiper';
export default {
  components: {
    Swiper,
    SwiperSlide,
  },
  setup() {
    const thumbsSwiper = ref(null);
    const setThumbsSwiper = (swiper) => {
        thumbsSwiper.value = swiper;
      };

    return {
      FreeMode,
      Thumbs,
      Navigation,
      thumbsSwiper,
      setThumbsSwiper,
    };
  },
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
</style>
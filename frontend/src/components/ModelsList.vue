<template>
  <div class="models-list">
    <transition-group name="models-list">
      <models-item
        v-for="model in models"
        :isAdmin="isAdmin"
        :userId="userId"
        :model="model"
        :key="model.model.id" 
        @remove="$emit('remove', model)"
      />
    </transition-group>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  props: { 
    models: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      isAdmin: false,
      userId: localStorage.getItem("id")
    }
  },
  mounted() {
    const accessToken = $cookies.get("access_token");
    const refreshToken = $cookies.get("refresh_token");

    const showUserInfo = (id) => {
      axios.get(`api/users/${id}`
      ).then(response => {
        this.isAdmin = response.data.is_admin;
        })
        .catch(error => {
          console.log(error);
        });
    };

    if (accessToken === null && refreshToken) {
      axios.post('/api/users/refresh', {
        refresh_token: refreshToken
      })
        .then(response => {
          $cookies.set('access_token', response.data.tokens.access_token, '15min', '/');
          $cookies.set('refresh_token', response.data.tokens.refresh_token, '7d', '/');
          localStorage.setItem('isAuth', true);
          showUserInfo(this.userId);
        })
        .catch(error => {
          console.log(error);
        });
    } else {
      showUserInfo(this.userId);
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
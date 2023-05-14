<template>
  <form class="dropzone">
    <div
      class="dropzone-container"
      @dragover="dragover"
      @dragleave="dragleave"
      @drop="drop"
    >
      <input
        type="file"
        multiple
        name="file"
        id="fileInput"
        class="hidden-input"
        @change="onChange"
        ref="file"
        accept="image/*, video/*, .glb, .gltf"
      />

      <label for="fileInput" class="file-label">
        <div v-if="isDragging">Отпустите, чтобы перетащить сюда файлы.</div>
        <div style="text-align: center;" v-else>Перетащите файлы сюда или нажмите здесь, чтобы загрузить.</div>
      </label>

      <div class="preview-container mt-4" v-if="files.length">
        <div v-for="file in files" :key="file.name" class="preview-card">
          <div>
            <img v-if="checkImages(file.name)" class="preview-img" :src="generateThumbnail(file)" />
            <p :title="file.name">
              {{ makeName(file.name) }}
            </p>
          </div>
          <div>
            <button
              class="ml-2"
              type="button"
              @click="remove(files.indexOf(file))"
              title="Удалить файл"
            >
              <b>&times;</b>
            </button>
          </div>
        </div>
      </div>
    </div>
  </form>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      isDragging: false,
      files: [],
      files_id: []
    };
  },
  methods: {
    checkImages(str) {
      return (/\.(jpeg|jpg|png)$/i).test(str);
    },
    onChange() {
      const accessToken = $cookies.get("access_token");
      const refreshToken = $cookies.get("refresh_token");
      this.files = [...this.$refs.file.files];
      this.$emit('onChange', {
        files_id: this.files_id
      })
      console.log(this.files);

      const saveFile = (file, access) => {
        let fileData = new FormData();  
        fileData.append('file', file);
        axios.post('/api/file',
          fileData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': `Bearer ${access}`
            }
          }
        ).then(response => {
            this.files_id.push(response.data.id);
            console.log(response);
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
            for (let i = 0; i < this.files.length; i++) {
              saveFile(this.files[i], accessToken);
            }
            console.log(response);
          })
          .catch(error => {
            console.log(error);
          });
      } else {
        for (let i = 0; i < this.files.length; i++) {
          saveFile(this.files[i], accessToken);
        }
      }
    },

    generateThumbnail(file) {
      let fileSrc = URL.createObjectURL(file);
      setTimeout(() => {
        URL.revokeObjectURL(fileSrc);
      }, 1000);
      return fileSrc;
    },

    makeName(name) {
      return (
        name.split(".")[0].substring(0, 3) +
        "..." +
        name.split(".")[name.split(".").length - 1]
      );
    },

    remove(i) {
      const accessToken = $cookies.get("access_token");
      const refreshToken = $cookies.get("refresh_token");
      const fileId = this.files_id[i];
      this.files.splice(i, 1);

      const removeFile = (id, access) => {
        axios.delete(`/api/file/${id}`,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': `Bearer ${access}`
            }
          }
        ).then(response => {
            const index = this.files_id.indexOf(id);
            this.files_id.splice(index, 1);
            console.log(response);
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
            removeFile(fileId, accessToken);
            console.log(response);
          })
          .catch(error => {
            console.log(error);
          });
      } else {
        removeFile(fileId, accessToken);
      }
    },

    dragover(e) {
      e.preventDefault();
      this.isDragging = true;
    },

    dragleave() {
      this.isDragging = false;
    },

    drop(e) {
      e.preventDefault();
      this.$refs.file.files = e.dataTransfer.files;
      this.onChange();
      this.isDragging = false;
    },
  },
};
</script>

<style>
.dropzone {
  margin-bottom: 1.5rem;
}
.dropzone-container {
  padding: 4rem;
  background: #f7fafc;
  border: 1px solid #e2e8f0;
}
.hidden-input {
  opacity: 0;
  overflow: hidden;
  position: absolute;
  width: 1px;
  height: 1px;
}
.file-label {
  font-size: 20px;
  display: block;
  cursor: pointer;
}
.preview-container {
  display: flex;
  flex-wrap: wrap;
  column-gap: 10px;
  row-gap: 10px;
  margin-top: 2rem;
}
.preview-card {
  display: flex;
  border: 1px solid #a2a2a2;
  padding: 5px;
  margin-left: 5px;
}
.preview-img {
  width: 50px;
  height: 50px;
  border-radius: 5px;
  border: 1px solid #a2a2a2;
  background-color: #a2a2a2;
}
</style>
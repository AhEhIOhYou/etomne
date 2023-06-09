<template>
  <form v-if="defaultFiles.length > 0" class="dropzone">
    <p>Загруженные файлы:</p>
    <div
      class="dropzone-container"
    >
      <input
        type="file"
        multiple
        name="defaultFile"
        id="defaultFileInput"
        class="hidden-input"
        ref="defaultFile"
        accept="image/*, video/*, .glb, .gltf"
      />

      <div v-if="defaultFiles.length" class="preview-container mt-4">
        <div v-for="file in defaultFiles" :key="file.title" class="preview-card">
          <button
            class="ml-2"
            type="button"
            @click="defaultRemove(defaultFiles.indexOf(file))"
            title="Удалить файл"
          >
            <b>&times;</b>
          </button>
          <div>
            <img v-if="checkImages(file.url)" class="preview-img" :src="'http://localhost:8095/' + file.url"/>
            <p :title="file.title">
              {{ file.title }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </form>
  <form class="dropzone">
    <p>Загружаемые файлы:</p>
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

      <div v-if="files.length" class="preview-container mt-4">
        <div v-for="file in files" :key="file.name" class="preview-card">
          <button
            class="ml-2"
            type="button"
            @click="defaultRemove(defaultFiles.indexOf(file))"
            title="Удалить файл"
          >
            <b>&times;</b>
          </button>
          <div>
            <img v-if="checkImages(file.name)" class="preview-img" :src="generateThumbnail(file)" />
            <p :title="file.name">
              {{ makeName(file.name) }}
            </p>
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
      defaultFiles: [],
      defaultFilesId: [],
      files: [],
      files_id: [],
    };
  },
  props: ['model'],
  mounted() {
    this.defaultFiles = this.getUrlsAndTitles(this.model.files);
    this.defaultFiles.forEach((file) => {
      this.defaultFilesId.push(file.id);
    });
  }, 
  methods: {
    getUrlsAndTitles(files) {
      const urlsAndTitles = [];
      for (const key in files) {
        if (files[key] && Array.isArray(files[key])) {
          files[key].forEach((file) => {
            urlsAndTitles.push({ title: file.title, url: file.url, id: file.id });
          });
        }
      }
      return urlsAndTitles;
    },
    checkImages(str) {
      return (/\.(jpeg|jpg|png)$/i).test(str);
    },
    onChange() {
      const accessToken = $cookies.get("access_token");
      const refreshToken = $cookies.get("refresh_token");
      const modelId = this.model.model.id;
      this.files = [...this.$refs.file.files];
      this.$emit('onChange', {
        files_id: this.files_id
      })

      const saveFile = (file, access) => {
        let fileData = new FormData();  
        fileData.append('file', file);
        
        axios.post(`api/model/addfile/${modelId}`,
          fileData,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': `Bearer ${access}`
            }
          }
        ).then(response => {
            this.files_id.push(response.data.id);
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
    defaultRemove(i) {
      const accessToken = $cookies.get("access_token");
      const refreshToken = $cookies.get("refresh_token");
      const fileId = this.defaultFilesId[i];
      this.defaultFiles.splice(i, 1);

      const removeFile = (id, access) => {
        axios.delete(`/api/file/${id}`,
          {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': `Bearer ${access}`
            }
          }
        ).then(response => {
            const index = this.defaultFilesId.indexOf(id);
            this.defaultFilesId.splice(index, 1);
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
          })
          .catch(error => {
            console.log(error);
          });
      } else {
        removeFile(fileId, accessToken);
      }
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

<style lang="scss">
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
  flex-direction: column;
  border: 1px solid #a2a2a2;
  padding: 5px;
  margin-left: 5px;

  & p {
    word-break: break-word;
  } 

  & .ml-2 {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0;
    margin-left: auto;
  }
}
.preview-img {
  width: 50px;
  height: 50px;
  border-radius: 5px;
  border: 1px solid #a2a2a2;
  background-color: #a2a2a2;
  object-fit: cover;
}
</style>
<template>
  <v-container>
      <v-flex xs12>
        <v-progress-linear :active="isLoading" class="ma-0" color="green lighten-3" height="4" indeterminate></v-progress-linear>
        <input id="fileUpload" type="file" accept="image/*" @change="previewFiles" hidden>
      </v-flex>
      <v-alert v-if="showAlert" type="error">
        Image size cannot exceed 2 MB
      </v-alert>
    <v-row class="ma-0">
        <v-img
             class="ma-1 selectable"
             max-height="109"
             max-width="109">
          <v-btn color="primary" id="uploadBtn" @click="chooseFiles()">
              <v-icon>mdi-camera</v-icon>
              Upload
          </v-btn>
        </v-img>
        <v-img
            v-for="(img, i) in images"
            :key="'img-' + i"
            :alt="img.alt"
            :src="img.src"
            class="ma-1 selectable"
            max-height="109"
            max-width="109"
            @click="selectImage(img)"
        >
            <v-col class="text-right ml-3">
                <v-btn v-if="i > 6" color="red" icon small @click.stop="deleteFile(i)">
                    <v-icon>mdi-trash-can-outline</v-icon>
                </v-btn>
            </v-col>
        </v-img>
     </v-row>
  </v-container>
</template>

<script>
    /**
     * Example of a custom Image selector
     * Key is to emit a select-file event when a file needs to be added
     */
    import { VImg } from 'vuetify/lib'
    import axios from "axios";
    import {default as API_ENDPOINTS} from "../api";
    import store from "../store/index"
    export default {
        name: "FileSelector",
        components: { VImg },
        data() {
            // Some public domain images from wikimedia.
            return {
                showAlert: false,
                isLoading: false,
                images: [
                    { src: 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Streifenhoernchen.jpg/1024px-Streifenhoernchen.jpg', alt: 'Siberian Chipmunk' },
                    { src: 'https://upload.wikimedia.org/wikipedia/commons/thumb/d/d8/NASA_Mars_Rover.jpg/750px-NASA_Mars_Rover.jpg', alt: 'NASA Mars Rover' },
                    { src: 'https://upload.wikimedia.org/wikipedia/commons/d/dd/Muybridge_race_horse_animated.gif', alt: 'Muybridge race horse animated' },
                    { src: 'https://upload.wikimedia.org/wikipedia/commons/2/2a/Locomotive_TEM2M-063_2006_G2.jpg', alt: 'Locomotive TEM2M-063 2006 G2' },
                    { src: 'https://upload.wikimedia.org/wikipedia/commons/8/80/ISS_March_2009.jpg', alt: 'ISS March 2009' },
                    { src: 'https://upload.wikimedia.org/wikipedia/commons/4/44/F-18F_after_launch_from_USS_Abraham_Lincoln_%28CVN-72%29.jpg', alt: 'F-18F after launch from USS Abraham Lincoln (CVN-72)' },
                    { src: 'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a8/Streifenhoernchen.jpg/1024px-Streifenhoernchen.jpg', alt: 'Siberian Chipmunk' }]
            }
        },
        methods: {
            selectImage(img) {
                this.$emit('select-file', img);
            },
            chooseFiles() {
                document.getElementById("fileUpload").click()
            },
            previewFiles(event) {
                const image = event?.target?.files[0]
                if (image) {
                    if (image.size / 1024 / 1024 >= 2) {
                        this.showAlert = true
                        setTimeout(() => {
                            this.showAlert = false
                        }, 4000)
                    } else {
                        this.uploadImage(image)
                    }
                }
            },
            uploadImage(image) {
                this.isLoading = true
                const formData = new FormData()
                formData.append("file", image)
                formData.append("requesterEmail", store.getters.email)
                return new Promise((resolve, reject) => {
                    axios({url: API_ENDPOINTS.UPLOAD_FILE, data: formData, method: 'POST' })
                        .then(resp => {
                            resolve(resp)
                            const newImage = {src: resp.data.url, alt: resp.data.url}
                            this.images.push(newImage)
                            this.isLoading = false
                        })
                        .catch(err => {
                            reject(err)
                            this.isLoading = false
                        })
                })
            },
            listImages() {
                return new Promise((resolve, reject) => {
                    axios({
                        url: API_ENDPOINTS.LIST_FILES(store.getters.email),
                        params: {requesterEmail: store.getters.email},
                        method: 'GET'
                    }).then(resp => {
                            resolve(resp)
                            resp.data.forEach(item => {
                                const newImage = {src: item.url, alt: ""}
                                this.images.push(newImage)
                            })
                            this.isLoading = false
                        }).catch(err => {
                            reject(err)
                            this.isLoading = false
                        })
                })
            },
            deleteFile(index) {
                this.isLoading = true
                const key = this.images[index].src?.split(store.getters.email + "/")[1]

                return new Promise((resolve, reject) => {
                    axios({
                        url: API_ENDPOINTS.DELETE_FILE(store.getters.email, key),
                        params: {requesterEmail: store.getters.email},
                        method: 'DELETE'
                    }).then(resp => {
                            resolve(resp)
                            this.images.splice(index, 1)
                            this.isLoading = false
                        }).catch(err => {
                            reject(err)
                            this.isLoading = false
                        })
                })
            }
        },
        created() {
            this.isLoading = true
            this.listImages()
        }
    }
</script>

<style scoped>
    .selectable {
        cursor: pointer;
    }
    #uploadBtn {
        width: 115px;
        height: 115px;
    }
</style>
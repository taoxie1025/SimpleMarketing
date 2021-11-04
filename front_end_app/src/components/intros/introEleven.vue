<template>
  <!-- intro-section -->
  <section :style="{ backgroundImage: 'url(' + projectBrief.backgroundImageUrl + ')' }"  id="intro-wrap" class="intro-eleven-Wrap text-white  text-center">
    <v-main class="bg-image">
      <b-container>
          <v-layout>
            <div id="avatar-wrap">
              <v-col class="ma-15">
                <v-btn class="mt-0 mr-2" fab large light :disabled="true" id="avatar-button">
                  <avatar v-if="projectBriefCopy.avatarUrl" alt="Avatar" :src="projectBriefCopy.avatarUrl" :size="60"></avatar>
                  <avatar v-else :username="projectBriefCopy.name" backgroundColor="grey" color="white"  :size="60"></avatar>
                </v-btn>
                <strong v-html="projectBriefCopy.author"></strong>
              </v-col>
            </div>
          </v-layout>
      </b-container>
      <b-container>
        <b-row>
          <b-col sm="12">
            <h1 class="intro-eleven-Title font-weight-bold text-42 t-shadow mb-3 text-white">
              {{projectBriefCopy.name}}
            </h1>
            <div class=" intro-eleven-Description text-16">
              <p class="mb-4">
                {{projectBriefCopy.intro}}
              </p>
            </div>

            <div class="intro-eleven-Buttons mb-5">
              <form onsubmit="return false"
                  class="d-flex align-items-center justify-content-center subscription_form m-auto"
              >
                <div class="form-group subscription-input-wrap">
                  <input
                      type="email"
                      required
                      class="form-control email_field"
                      v-model="subscribeForm.email"
                      aria-describedby="emailHelpId"
                      placeholder="Please enter your email"
                  />
                  <button @click="subscribe" type="submit" class="btn btn-submit btn-gradient">
                    Subscribe
                  </button>
                </div>
              </form>
              <div id="message" v-if="showMessage"
                   class="intro-eleven-Description alert alert-dismissible fade show"
                   role="alert"
              >
                <button type="button" class="close" @click="showMessage=false">
                  <span aria-hidden="true">&times;</span>
                </button>
                <strong>{{displayedMessage}}</strong>
              </div>
              <div class="form-group"></div>
            </div>
          </b-col>
        </b-row>
      </b-container>
    </v-main>
    <div class="overlay"></div>
  </section>

  <!-- end::intro-section -->
</template>

<script>
import Avatar from 'vue-avatar'
import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  metaInfo: {
    bodyAttrs: {
      class: ["landing-gradient-steel-gray"]
    }
  },
  components: {
    Avatar
  },
  data () {
    return {
      projectBrief: {
        projectId: "",
        name: "",
        createdAt: "",
        interval: "",
        lastBroadcastTimeMs: "",
        intro: "",
        backgroundImageUrl: "",
        avatarUrl: "",
        outgoingEmail: "",
        author: "",
      },
      projectBriefCopy: {},
      subscribeForm: {
        email:"",
        firstName: "",
        lastName: ""
      },
      showMessage: false,
      displayedMessage: ""
    }
  },
  methods: {
    getProjectBrief(projectId) {
      return new Promise((resolve) => {
        axios({url: API_ENDPOINTS.READ_PROJECT_BRIEF(projectId), method: 'GET'})
            .then(resp => {
              resolve(resp)
              this.projectBrief = resp.data
              this.projectBriefCopy = this.projectBrief
            })
            .catch(() => {
              //TODO: directs to 404 page
            })
      })
    },
    isEmailValid(email) {
      if (!(/.+@.+\..+/.test(email))) {
        return false
      }
      return true
    },
    subscribe() {
      if (!this.isEmailValid(this.subscribeForm.email)) {
        return
      }
      const data = {projectId: this.projectBrief.projectId, email: this.subscribeForm.email}
      return new Promise((resolve) => {
        axios({url: API_ENDPOINTS.CREATE_SUBSCRIBER, data: data, method: 'POST'})
            .then(resp => {
              resolve(resp)
              this.displayedMessage = "Subscribed successfully!"
              this.showMessage = true
              this.subscribeForm.email = ""
              this.closeShowMessage()
            })
            .catch(() => {
              this.displayedMessage = "Failed to subscribe."
              this.showMessage = true
              this.closeShowMessage()
            })
      })
    },
    closeShowMessage() {
      setTimeout(() => {
        this.showMessage = false
      }, 6000)
    },
  },
  created() {
    const projectId = "project-" + this.$route.params.projectIdShort
    this.getProjectBrief(projectId)
  }
};
</script>

<style>
.bg-image {
  background-repeat: no-repeat;
  background-position: center center;
  background-attachment: fixed;
  background-size: cover;
  min-height: 60vh;
}
#avatar-wrap {
  margin-left: -150px;
  margin-bottom: 180px;
  margin-top: -150px;
}
#message {
  width: 350px;
}
</style>
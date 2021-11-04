<template>
  <v-app id="app">
    <v-container fluid>
      <v-row justify="center" align="center">
        <v-col cols="12" lg="8">
          <h2>Unsubscribe</h2>
        </v-col>
        <v-col cols="12" lg="8">
          <v-card>
            <v-card-text>
              <span>Enter your email address below to unsubscribe</span>
              <span v-if="projectName">&nbsp;from {{ projectName}}</span>
            </v-card-text>
            <v-form ref="unsubscribeForm" v-model="isUnsubscribeFormValid" lazy-validation>
              <v-text-field class="ml-10 mr-10" v-model="email" :rules="rules" :disabled="isUpdating || !projectId" color="blue-grey lighten-2" label="Email" outlined required dense></v-text-field>
              <v-btn class="ml-10 mb-3" tile :disabled="!isUnsubscribeFormValid || isUpdating" @click="unsubscribe">Unsubscribe</v-btn>
            </v-form>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-snackbar v-model="snackbar" timeout="3000" @click="snackbar=false">
          {{quickMessage}}
          <template v-slot:action="{ attrs }">
            <v-btn
                color="pink"
                text
                v-bind="attrs"
                @click="snackbar = false"
            >
              Close
            </v-btn>
          </template>
        </v-snackbar>
      </v-row>
    </v-container>
  </v-app>
</template>

<script>

import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  components: {},
  data() {
    return {
      projectId: "",
      projectName: "",
      email: "",
      isUnsubscribeFormValid: false,
      isUpdating: false,
      snackbar: false,
      quickMessage: "",
      rules: [
        v => !!v || "Required",
        v => /.+@.+\..+/.test(v) || "E-mail must be valid"
      ],
    }
  },
  watch: {
  },
  methods: {
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
    },
    reset() {
      this.email = ""
    },
    unsubscribe() {
      if (this.$refs.unsubscribeForm.validate()) {
        return new Promise((resolve, reject) => {
          this.isUpdating = true
          axios({url: API_ENDPOINTS.DELETE_SUBSCRIBER(this.projectId, this.email), method: 'DELETE'})
              .then(resp => {
                resolve(resp)
                this.isUpdating = false
                this.reset()
                this.$refs.unsubscribeForm.reset();
                this.showQuickMessage("Unsubscribe successfully.")
              })
              .catch(err => {
                reject(err)
                this.isUpdating = false
                this.showQuickMessage("Failed to unsubscribe")
              })
        })
      }
    }
  },
  computed: {
  },
  mounted() {
  },
  created() {
    this.projectId = this.$route.query.projectId
    this.projectName = this.$route.query.projectName
  }
}
</script>


<style>

</style>

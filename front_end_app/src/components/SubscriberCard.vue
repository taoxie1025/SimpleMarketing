<template>
  <v-card flat class="mx-auto">
    <v-card-title class="headline grey lighten-2">
      <v-icon left color="black">mdi-file-document-edit-outline</v-icon>
      <h3 v-if="!subscriberCopy.email">New Subscriber</h3>
      <h3 v-else>Edit Subscriber</h3>
    </v-card-title>
    <v-container>
      <br>
      <v-flex xs12 sm12>
        <v-text-field
            label="Email Address"
            placeholder="New subscriber"
            outlined
            :prepend-inner-icon="!isEditing ? 'mdi-pencil-plus-outline' : ''"
            v-model="subscriberCopy.email"
            :disabled="isEditing"
            :rules="emailRules"
            block
            required
        >
        </v-text-field>
      </v-flex>
      <v-row xs12 class="mx-auto" v-if="isEditing">
        <v-flex xs6 sm6 pr-1>
          <v-text-field
              label="Current Article"
              placeholder="Article No."
              outlined
              prepend-inner-icon="mdi-pencil-plus-outline"
              v-model.number="subscriberCopy.articleCursor"
              :rules="articleCursorRules"
              block
              required
          >
          </v-text-field>
        </v-flex>
        <v-flex xs6 sm6 pl-1>
          <v-select
              v-model="subscriberCopy.isEnabled"
              :items="subscriberStatus"
              item-text="option"
              item-value="value"
              label="Status"
              block
              outlined
          ></v-select>
        </v-flex>
      </v-row>
      <v-flex xs12 sm12>
        <div id="message" v-if="showMessage"
             class="intro-eleven-Description alert alert-dismissible fade show small success"
             role="alert"
        >
          <button type="button" class="close" @click="showMessage=false">
            <span aria-hidden="true">&times;</span>
          </button>
          <strong>{{displayedMessage}}</strong>
        </div>
      </v-flex>
      <v-row class="mx-auto">
        <v-btn class="ma-1" tile outlined color="info" @click="close">Cancel</v-btn>
        <v-spacer></v-spacer>
        <v-btn v-if="!isEditing" class="ma-1" tile color="info" @click="add" :disabled="!isEmailValid(subscriberCopy.email)">Add</v-btn>
        <v-btn v-else class="ma-1" dark tile color="info" @click="save">Save</v-btn>
      </v-row>
    </v-container>
  </v-card>
</template>

<script>


import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  name: "SubscriberCard",
  components: {},
  props: [
      'projectId',
      'email',
      'subscriber',
      'isEditing'
  ],
  data() {
    return {
      subscriberCopy: {},
      emailRules: [
        v => !!v || "Required",
        v => /.+@.+\..+/.test(v) || "E-mail must be valid"
      ],
      articleCursorRules: [
        v => !!v || "Required",
        v => /^\d+$/.test(v) || "Article No. must be a number"
      ],
      showMessage: false,
      displayedMessage: "",
      subscriberStatus: [
        {option: "Enabled", value: true},
        {option: "Disabled", value: false}
      ]
    }
  },
  computed: {
  },
  methods: {
    close() {
      this.$emit('close')
    },
    closeShowMessage() {
      setTimeout(() => {
        this.showMessage = false
      }, 6000)
    },
    isEmailValid(email) {
      if (!(/.+@.+\..+/.test(email))) {
        return false
      }
      return true
    },
    isArticleCursorValid(articleNumber) {
      if (!(/^\d+$/.test(articleNumber))) {
        return false
      }
      return true
    },
    save() {
      if (!this.isArticleCursorValid(this.subscriberCopy.articleCursor)) {
        return
      }
      this.saveSubscriber()
    },
    add() {
      if (!this.isEmailValid(this.subscriberCopy.email)) {
        return
      }
      this.subscribe()
    },
    subscribe() {
      const data = {projectId: this.projectId, email: this.subscriberCopy.email, firstName: '', lastName: ''}
      return new Promise((resolve) => {
        axios({url: API_ENDPOINTS.CREATE_SUBSCRIBER, data: data, method: 'POST'})
            .then(resp => {
              resolve(resp)
              this.displayedMessage = "Subscribed successfully!"
              this.showMessage = true
              this.subscriberCopy.email = ""
              this.closeShowMessage()
              this.$emit('subscriberAdded', (resp.data))
            })
            .catch(() => {
              this.displayedMessage = "Failed to subscribe."
              this.showMessage = true
              this.closeShowMessage()
            })
      })
    },
    saveSubscriber() {
      const data = {projectId: this.projectId, email: this.subscriberCopy.email, firstName: '', lastName: '',
        isEnabled: this.subscriberCopy.isEnabled, articleCursor: this.subscriberCopy.articleCursor}
      return new Promise((resolve) => {
        axios({url: API_ENDPOINTS.UPDATE_SUBSCRIBER, data: data, method: 'PUT'})
            .then(resp => {
              resolve(resp)
              Object.assign(this.subscriber, resp.data)
              this.$emit('subscriberUpdated', resp.data)
            })
            .catch(() => {
            })
      })
    }
  },
  created() {
    Object.assign(this.subscriberCopy, this.subscriber)
    if (this.subscriberCopy.isEnabled == "Yes" || this.subscriber.isEnabled) {
      this.subscriberCopy.isEnabled = true
    } else {
      this.subscriberCopy.isEnabled = false
    }
  }
}
</script>

<style scoped>

</style>
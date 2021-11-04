<template>
  <v-app id="app">
    <v-container fluid>
      <v-layout row wrap class="ma-1">
        <v-flex xs12>
          <h3>Subscribers</h3>
        </v-flex>
        <v-flex xs12>
            <v-text-field
                label="Email Address"
                placeholder="Search for subscriber"
                outlined
                prepend-inner-icon="mdi-magnify"
                :append-icon="emailFilter != '' ? 'mdi-close' : ''"
                v-model="emailFilter"
                @click:append="clear"
            >
            </v-text-field>
        </v-flex>
        <v-flex xs12 class="mb-3">
            <v-btn tile dark color="info" @click="showAddSubscriberDialog">
              <v-icon left>mdi-plus-thick</v-icon>
              Add subscriber
            </v-btn>
            <v-btn class="ma-1" tile outlined color="info" :disabled="selectedSubscribers.length==0" @click="showSubscriberDialogOnEdit">Edit</v-btn>
            <v-btn class="ma-1" tile outlined color="info" :disabled="selectedSubscribers.length==0" @click="showConfirmDialog">Delete</v-btn>
        </v-flex>
        <v-flex xs12>
          <v-data-table
              @click:row="rowClick"
              v-model="selectedSubscribers"
              :headers="headers"
              :items="subscribers"
              item-key="email"
              :loading="isLoading"
              loading-text="Loading..."
              hide-default-footer
              :server-items-length="pageSize"
              disable-sort
              single-select
              class="elevation-0"
              v-infinite-scroll="loadMore"
              infinite-scroll-disabled="busy"
              infinite-scroll-distance="10"
              infinite-scroll-immediate-check="false"
          >
          </v-data-table>
        </v-flex>
       </v-layout>
      <v-layout>
        <v-flex justify="center">
          <v-dialog v-if="subscriberEditorDialog" v-model="subscriberEditorDialog" width="500px">
              <subscriber-form
                  :projectId="$route.params.projectId"
                  :subscriber="selectedSubscribers.length>0 ? selectedSubscribers[0] : {}"
                  :isEditing="selectedSubscribers.length>0 ? true : false"
                  @close="subscriberEditorDialog=false"
                  @subscriberAdded="onSubscriberAdded"
                  @subscriberUpdated="onSubscriberUpdated"></subscriber-form>
          </v-dialog>
        </v-flex>
      </v-layout>
      <v-layout>
        <v-row justify="center">
          <v-dialog v-model="deleteConfirmDialog" max-width="250">
            <confirm-dialog v-bind:title="`Are you sure?`" :body="`The subscriber will no longer receive your content.`" @yes="deleteConfirmed" @no="deleteConfirmDialog=false"></confirm-dialog>
          </v-dialog>
        </v-row>
      </v-layout>
    </v-container>
  </v-app>
</template>

<script>

import axios from "axios";
import {default as API_ENDPOINTS} from "../api";
import SubscriberForm from "@/components/SubscriberCard";
import ConfirmDialog from "@/components/ConfirmDialog";

export default {
  components: {SubscriberForm, ConfirmDialog},
  data() {
    return {
      headers: [
        { text: 'Email', align: 'start', value: 'email', class: "grey lighten-4"},
        { text: 'Current Article No.', align: 'center', value: 'articleCursor', class: "grey lighten-4"},
        { text: 'Last Broadcast Date', align: 'center', value: 'lastBroadcastTimeMs', class: "grey lighten-4"},
        { text: 'Subscription Date', align: 'center', value: 'createdAt', class: "grey lighten-4"},
        { text: 'Enabled', align: 'end', value: 'isEnabled', class: "grey lighten-4"}
        ],
      token: "",
      pageSize: 15,
      subscribers: [],
      isLoading: false,
      selectedSubscribers: [],
      emailFilter: "",
      subscriberEditorDialog: false,
      deleteConfirmDialog: false
    }
  },
  watch: {
    emailFilter() {
      if (this.emailFilter == "") {
        this.subscribers = []
        this.readSubscribers()
      } else {
        this.subscribers = []
        this.token = ""
        this.searchSubscriber()
      }
    }
  },
  methods: {
    onSubscriberUpdated(updatedSubscriber) {
      this.subscriberEditorDialog = false
      let index = 0;
      for (let i = 0; i < this.subscribers.length; i++) {
        if (this.subscribers[i].email == updatedSubscriber.email) {
          index = i
          break
        }
      }
      if (updatedSubscriber.isEnabled) {
        updatedSubscriber.isEnabled = "Yes"
      } else {
        updatedSubscriber.isEnabled = "No"
      }
      updatedSubscriber.createdAt = this.convertEpochToDateString(updatedSubscriber.createdAt)
      updatedSubscriber.lastBroadcastTimeMs = this.convertEpochToDateString(updatedSubscriber.lastBroadcastTimeMs)
      this.subscribers.splice(index, 1, updatedSubscriber)
    },
    onSubscriberAdded(newSubscriber) {
      this.subscribers.unshift(newSubscriber)
    },
    deleteConfirmed() {
      let index = 0;
      for (let i = 0; i < this.subscribers.length; i++) {
        if (this.subscribers[i].email == this.selectedSubscribers[0].email) {
          index = i
          break
        }
      }
      return new Promise((resolve, reject) => {
        axios({url: API_ENDPOINTS.DELETE_SUBSCRIBER(this.$route.params.projectId, this.selectedSubscribers[0].email), method: 'DELETE' })
            .then(resp => {
              resolve(resp)
              this.subscribers.splice(index, 1)
              this.deleteConfirmDialog = false
            })
            .catch(err => {
              reject(err)
              this.deleteConfirmDialog = false
            })
      })
    },
    showConfirmDialog() {
      this.deleteConfirmDialog = true
    },
    rowClick(item) {
      this.selectedSubscribers = []
      this.selectedSubscribers.push(item)
    },
    highlightClickedRow(item) {
      const tr = item.target.parentNode;
      tr.classList.add('highlight');
    },
    showAddSubscriberDialog() {
      this.selectedSubscribers = []
      this.subscriberEditorDialog = true
    },
    showSubscriberDialogOnEdit() {
      this.subscriberEditorDialog = true
    },
    clear() {
      this.emailFilter = ""
    },
    searchSubscriber() {
      this.isLoading = true
      this.subscribers = []
      this.selectedSubscribers = []
      return new Promise((resolve, reject) => {
        let data = {projectId: this.$route.params.projectId, emailFilter: this.emailFilter}
        axios({url: API_ENDPOINTS.SEARCH_SUBSCRIBERS, params: data, method: 'GET' })
            .then(resp => {
              resolve(resp)
              for (let i = 0; i < resp.data.length; i++) {
                const subscriber = {email: resp.data[i].email,
                  articleCursor: resp.data[i].articleCursor,
                  lastBroadcastTimeMs: this.convertEpochToDateString(resp.data[i].lastBroadcastTimeMs),
                  createdAt: this.convertEpochToDateString(resp.data[i].createdAt),
                  isEnabled: resp.data[i].isEnabled ? "Yes" : "No"
                }
                this.subscribers.push(subscriber)
              }
              this.isLoading = false
            })
            .catch(err => {
              reject(err)
              this.isLoading = false
            })
      })
    },
    loadMore() {
      if (this.token != "EOF" && this.emailFilter == "") {
        this.readSubscribers()
      }
    },
    validateLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        console.log("not signed in ");
        this.$router.push('/')
        return false
      }
      return true
    },
    readSubscribers() {
      this.isLoading = true
      return new Promise((resolve, reject) => {
        let data = {email: this.$store.getters.email ? this.$store.getters.email : this.project.email,
          projectId: this.$route.params.projectId, pageSize: this.pageSize,
          token: this.token
        }
        axios({url: API_ENDPOINTS.READ_SUBSCRIBERS, data: data, method: 'POST' })
            .then(resp => {
              resolve(resp)
              this.token = resp.data.token
              for (let i = 0; i < resp.data.subscribers.length; i++) {
                const subscriber = {email: resp.data.subscribers[i].email,
                  articleCursor: resp.data.subscribers[i].articleCursor,
                  lastBroadcastTimeMs: this.convertEpochToDateString(resp.data.subscribers[i].lastBroadcastTimeMs),
                  createdAt: this.convertEpochToDateString(resp.data.subscribers[i].createdAt),
                  isEnabled: resp.data.subscribers[i].isEnabled ? "Yes" : "No"
                }
                this.subscribers.push(subscriber)
              }
              this.isLoading = false
            })
            .catch(err => {
              reject(err)
              this.isLoading = false
            })
      })
    },
    convertEpochToDateString(timestamp) {
      if (timestamp && timestamp > 0) {
        const date = new Date(timestamp)
        return date.toLocaleDateString()
      }
      return "-"
    }
  },
  created: function () {
    this.validateLogin()

    //bypass javascript not enabled issue
    if (!this.subscribers || this.subscribers.length == 0) {
      this.subscribers = []
      this.loadMore()
    }
  }
}
</script>


<style>
.v-text-field .v-icon {
  color: grey !important;
}
.v-data-table .v-icon {
  color: grey !important;
}
</style>

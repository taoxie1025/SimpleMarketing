<template>
  <v-app id="app">
    <v-container fluid>
      <v-layout row wrap class="ma-1">
        <v-flex xs12>
          <v-text-field
              label="Email Address"
              placeholder="Search for subscriber"
              outlined
              prepend-inner-icon="mdi-magnify"
              :append-icon="subscriberFilter != '' ? 'mdi-close' : ''"
              v-model="subscriberFilter"
              @click:append="clear"
          >
          </v-text-field>
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
            <admin-subscriber-form
                :projectId="selectedSubscribers[0].projectId"
                :subscriber="selectedSubscribers.length>0 ? selectedSubscribers[0] : {}"
                :isEditing="selectedSubscribers.length>0 ? true : false"
                @close="subscriberEditorDialog=false"
                @subscriberUpdated="onSubscriberUpdated">
            </admin-subscriber-form>
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
import AdminSubscriberForm from "@/components/AdminSubscriberCard";
import ConfirmDialog from "@/components/ConfirmDialog";

export default {
  components: {AdminSubscriberForm, ConfirmDialog},
  data() {
    return {
      headers: [
        { text: 'Email', align: 'start', value: 'email', class: "grey lighten-4"},
        { text: 'ProjectId', align: 'start', value: 'projectId', class: "grey lighten-4"},
        { text: 'Current Article No.', align: 'center', value: 'articleCursor', class: "grey lighten-4"},
        { text: 'Last Broadcast Date', align: 'center', value: 'lastBroadcastTimeMsLabel', class: "grey lighten-4"},
        { text: 'Subscription Date', align: 'center', value: 'createdAtLabel', class: "grey lighten-4"},
        { text: 'Enabled', align: 'end', value: 'isEnabledLabel', class: "grey lighten-4"}
      ],
      token: "",
      pageSize: 15,
      subscribers: [],
      isLoading: false,
      selectedSubscribers: [],
      subscriberFilter: "",
      subscriberEditorDialog: false,
      deleteConfirmDialog: false
    }
  },
  watch: {
    subscriberFilter() {
      if (this.subscriberFilter == "") {
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
        updatedSubscriber.isEnabledLabel = "Yes"
      } else {
        updatedSubscriber.isEnabledLabel = "No"
      }
      updatedSubscriber.createdAtLabel = this.convertEpochToDateString(this.subscribers[index].createdAt)
      updatedSubscriber.lastBroadcastTimeMsLabel = this.convertEpochToDateString(this.subscribers[index].lastBroadcastTimeMs)
      this.subscribers.splice(index, 1, updatedSubscriber)
    },
    rowClick(item) {
      this.selectedSubscribers = []
      this.selectedSubscribers.push(item)
      this.subscriberEditorDialog = true
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
      this.subscriberFilter = ""
    },
    searchSubscriber() {
      this.isLoading = true
      this.subscribers = []
      this.selectedSubscribers = []
      return new Promise((resolve, reject) => {
        let data = {projectId: this.$route.params.projectId, subscriberFilter: this.subscriberFilter}
        axios({url: API_ENDPOINTS.ADMIN_SEARCH_SUBSCRIBERS, params: data, method: 'GET' })
            .then(resp => {
              resolve(resp)
              for (let i = 0; i < resp.data.length; i++) {
                const subscriber = {
                  email: resp.data[i].email,
                  projectId: resp.data[i].projectId,
                  articleCursor: resp.data[i].articleCursor,
                  lastBroadcastTimeMs: this.convertEpochToDateString(resp.data[i].lastBroadcastTimeMs),
                  lastBroadcastTimeMsLabel: resp.data[i].lastBroadcastTimeMs,
                  createdAtLabel: this.convertEpochToDateString(resp.data[i].createdAt),
                  createdAt: resp.data[i].createdAt,
                  isEnabled: resp.data[i].isEnabled,
                  isEnabledLabel: resp.data[i].isEnabled ? "Yes" : "No"
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
      if (this.token != "EOF" && this.subscriberFilter == "") {
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
        axios({url: API_ENDPOINTS.ADMIN_READ_SUBSCRIBERS, data: data, method: 'POST' })
            .then(resp => {
              resolve(resp)
              this.token = resp.data.token
              for (let i = 0; i < resp.data.subscribers.length; i++) {
                const subscriber = {
                  email: resp.data.subscribers[i].email,
                  projectId: resp.data.subscribers[i].projectId,
                  articleCursor: resp.data.subscribers[i].articleCursor,
                  lastBroadcastTimeMs: resp.data.subscribers[i].lastBroadcastTimeMs,
                  lastBroadcastTimeMsLabel: this.convertEpochToDateString(resp.data.subscribers[i].lastBroadcastTimeMs),
                  createdAt: resp.data.subscribers[i].createdAt,
                  createdAtLabel: this.convertEpochToDateString(resp.data.subscribers[i].createdAt),
                  isEnabled: resp.data.subscribers[i].isEnabled,
                  isEnabledLabel: resp.data.subscribers[i].isEnabled ? "Yes" : "No"
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

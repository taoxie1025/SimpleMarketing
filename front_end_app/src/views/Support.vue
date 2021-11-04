<template>
  <v-app id="app">
    <v-container fluid>
      <v-layout row wrap class="ma-1">
        <v-flex xs12 class="mb-4">
          <div id="button-wrap">
            <v-btn tile dark color="info" @click="createNewTicket">
              <v-icon left>mdi-plus</v-icon>
              Create New Ticket
            </v-btn>
          </div>
        </v-flex>
        <v-flex xs12>
          <v-expand-transition>
            <ticket-card v-if="isCreateNewTicket" v-model="isCreateNewTicket" isNewTicket="true" :ticket="newTicket" :projects="projectLabels" @cancel="onCancel" @createTicket="onCreateTicket"></ticket-card>
          </v-expand-transition>
        </v-flex>
        <v-flex xs12 class="mt-5">
          <h2>Your tickets</h2>
        </v-flex>
        <v-flex xs12>
          <ticket-table :tickets="tickets" @appendComment="onAppendComment" v-on:click.native="ticketClicked"></ticket-table>
        </v-flex>
        <v-flex xs12>
          <div v-if="token != 'EOF'">
          <v-layout justify-center>
            <v-card-actions>
              <v-btn tile small depressed style="background-color: white; color: #5bc0de" @click="readTickets">LOAD MORE</v-btn>
            </v-card-actions>
          </v-layout>
          </div>
        </v-flex>
        <v-snackbar v-model="snackbar" timeout="8000" @click="snackbar=false">
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
      </v-layout>
    </v-container>
  </v-app>
</template>

<script>

import {EventBus as bus} from "../event_bus.js"
import TicketTable from "@/components/TicketTable";
import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";
import TicketCard from "@/components/TicketCard";

export default {
  components: {TicketCard, TicketTable},
  data() {
    return {
      projects: [], // to provide a selection when creating new ticket
      projectLabels: [], // simple key value paris to parse down to ticket card for selections
      isLoading: false,
      tickets: [],
      newTicket: {},
      isUpdating: false,
      isCreateNewTicket: false,
      snackbar: false,
      quickMessage: "",
      token: "",
      pageSize: 10,
    }
  },
  watch: {
  },
  methods: {
    ticketClicked() {
      this.isCreateNewTicket=false
    },
    onAppendComment(comment) {
      return new Promise((resolve, reject) => {
        this.isUpdating = true
        const data = comment
        data.email = this.$store.getters.email
        data.name = this.$store.getters.email
        axios({url: API_ENDPOINTS.CREATE_COMMENT, data: data, method: 'POST'})
            .then(resp => {
              resolve(resp)
              for (let i = 0; i < this.tickets.length; i++) {
                if (this.tickets[i].ticketId == resp.data.ticketId) {
                  this.tickets[i].comments.push(resp.data)
                  break
                }
              }
              bus.$emit('commentAppended')
              this.isUpdating = false
            })
            .catch(err => {
              reject(err)
              this.isUpdating = false
              this.showQuickMessage("Failed to reply the ticket, please try again later.")
            })
      })
    },
    onCreateTicket(projectId) {
      console.log("XXXXXXX " + projectId + " " + this.getProjectName(projectId))
      return new Promise((resolve, reject) => {
        this.isUpdating = true
        const data = this.newTicket
        data.email = this.$store.getters.email
        data.projectName = this.getProjectName(projectId)
        axios({url: API_ENDPOINTS.CREATE_TICKET, data: data, method: 'POST'})
            .then(resp => {
              resolve(resp)
              this.isUpdating = false
              this.tickets.splice(0, 0, resp.data)
              this.onCancel()
              this.newTicket = {}
              this.showQuickMessage("Your ticket is created successfully.")
            })
            .catch(err => {
              reject(err)
              this.isUpdating = false
              this.showQuickMessage("Failed to create ticket.")
            })
      })
    },
    onCancel() {
      this.isCreateNewTicket = false
    },
    createNewTicket(){
      this.isCreateNewTicket = true
      bus.$emit('creatingNewTicket')
    },
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
    },
    readTickets() {
      return new Promise((resolve, reject) => {
        this.isUpdating = true
        const data = {
          email: this.$store.getters.email,
          token: this.token,
          pageSize: this.pageSize
        }
        axios({url: API_ENDPOINTS.READ_TICKETS, data: data, method: 'POST'})
            .then(resp => {
              resolve(resp)
              this.isUpdating = false
              for (let i = 0; i < resp.data.tickets.length; i++) {
                this.tickets.push(resp.data.tickets[i])
              }
              this.token = resp.data.token
            })
            .catch(err => {
              reject(err)
              this.isUpdating = false
              this.showQuickMessage("Failed to fetch tickets, please try again.")
            })
      })
    },
    readUserAndProjects() {
      this.$store.dispatch('readUser', this.$store.getters.email)
          .then(() => {
            if (this.$store.getters.user.projectIds && this.$store.getters.user.projectIds.length > 0) {
              this.readProjects(this.$store.getters.email, this.$store.getters.user.projectIds).then(resp => {
                const sortedProjects = resp.data.sort(function(x, y) {
                  return x.createdAt < y.createdAt ? 1 : (x.createdAt > y.createdAt ? -1 : 0)
                })
                this.projects = sortedProjects
                this.convertProjects()
                this.isLoading = false
              })
            } else {
              this.isLoading = false
            }
          })
    },
    readProjects(email, projectIds){
      const data = {email, projectIds}
      return new Promise((resolve, reject) => {
        axios({url: API_ENDPOINTS.READ_PROJECTS, data: data, method: 'POST' })
            .then(resp => {
              resolve(resp)
            })
            .catch(err => {
              reject(err)
            })
      })
    },
    convertProjects() {
      for (let i = 0; i < this.projects.length; i++) {
        this.projectLabels.push({option: this.projects[i].name, value: this.projects[i].projectId})
      }
    },
    getProjectName(projectId) {
      for (let i = 0; i < this.projects.length; i++) {
        if (this.projects[i].projectId == projectId) {
          return this.projects[i].name
        }
      }
      return ""
    },
    validateLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        this.$router.push('/')
        return false
      }
      return true
    }
  },
  mounted() {
    bus.$emit('commentAppended')
  },
  created: function () {
    if (this.validateLogin()) {
      this.readTickets()
      this.readUserAndProjects()
    }
  }
}
</script>


<style>

</style>

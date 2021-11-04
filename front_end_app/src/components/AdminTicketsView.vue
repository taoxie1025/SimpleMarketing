<template>
  <v-app id="app">
    <v-container fluid>
      <v-layout row wrap class="ma-1">
        <v-flex xs12>
          <v-text-field
              label="Ticket"
              placeholder="Search for ticket"
              outlined
              prepend-inner-icon="mdi-magnify"
              :append-icon="ticketFilter != '' ? 'mdi-close' : ''"
              v-model="ticketFilter"
              @click:append="clear"
          >
          </v-text-field>
        </v-flex>
        <v-flex xs12>
          <v-data-table
              @click:row="rowClick"
              v-model="selectedTickets"
              :headers="headers"
              :items="tickets"
              item-key="ticketId"
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
          <v-dialog v-if="ticketBriefDialog" v-model="ticketBriefDialog" width="1000px">
            <v-card tile>
              <v-expansion-panels id="panel" v-model="panels">
                <v-expansion-panel>
                  <ticket-brief
                      :ticket="selectedTickets[0]"
                      @appendComment="onAppendComment"
                      @cancel="onCancel"
                      @ticketUpdated="onTicketUpdated"
                      @ticketDeleted="onTicketDeleted">
                  >
                  </ticket-brief>
                </v-expansion-panel>
              </v-expansion-panels>
            </v-card>
          </v-dialog>
        </v-flex>
      </v-layout>
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
    </v-container>
  </v-app>
</template>

<script>

import axios from "axios";
import {default as API_ENDPOINTS} from "../api";
import TicketBrief from "@/components/TicketBrief";

export default {
  components: {TicketBrief},
  data() {
    return {
      headers: [
        { text: 'Type', align: 'start', value: 'ticketTypeLabel', class: "grey lighten-4"},
        { text: 'User Email', align: 'start', value: 'email', class: "grey lighten-4"},
        { text: 'Status', align: 'start', value: 'ticketStatusLabel', class: "grey lighten-4"},
        { text: 'Replies', align: 'start', value: 'reply', class: "grey lighten-4"},
        { text: 'Date Created', align: 'start', value: 'createdAtLabel', class: "grey lighten-4"},
        { text: 'Title', align: 'start', value: 'title', class: "grey lighten-4"},
        { text: 'Project ID', align: 'start', value: 'projectId', class: "grey lighten-4"},
        { text: 'Project Name', align: 'start', value: 'projectName', class: "grey lighten-4"}
      ],
      token: "",
      pageSize: 15,
      tickets: [],
      isLoading: false,
      selectedTickets: [],
      ticketFilter: "",
      ticketBriefDialog: false,
      panels: 0,
      isUpdating: false,
      quickMessage: "",
      snackbar: false
    }
  },
  watch: {
    ticketFilter() {
      if (this.ticketFilter == "") {
        this.tickets = []
        this.readTickets()
      } else {
        this.tickets = []
        this.token = ""
        this.searchTickets()
      }
    }
  },
  methods: {
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
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
                  this.tickets[i].reply += 1
                  break
                }
              }
              this.isUpdating = false
            })
            .catch(err => {
              reject(err)
              this.isUpdating = false
              this.showQuickMessage("Failed to reply the ticket, please try again later.")
            })
      })
    },
    onTicketUpdated(updatedTicket) {
      for (let i = 0; i < this.tickets.length; i++) {
        if (this.tickets[i].ticketId == updatedTicket.ticketId) {
          this.tickets.splice(i, 1, updatedTicket)
          this.ticketBriefDialog = false
          if (this.tickets[i].ticketStatus != 1) {
            this.tickets.splice(i, 1)
          }
          break
        }
      }
    },
    onTicketDeleted(ticketId) {
      for (let i = 0; i < this.tickets.length; i++) {
        if (this.tickets[i].ticketId == ticketId) {
          this.tickets.splice(i, 1)
          this.ticketBriefDialog = false
          break
        }
      }
    },
    onCancel() {
      this.ticketBriefDialog = false
    },
    rowClick(item) {
      this.selectedTickets = []
      this.selectedTickets.push(item)
      this.ticketBriefDialog = true
    },
    highlightClickedRow(item) {
      const tr = item.target.parentNode;
      tr.classList.add('highlight');
    },
    clear() {
      this.ticketFilter = ""
    },
    searchTickets() {
      this.isLoading = true
      this.tickets = []
      this.selectedTickets = []
      return new Promise((resolve, reject) => {
        let data = {ticketFilter: this.ticketFilter}
        axios({url: API_ENDPOINTS.ADMIN_SEARCH_TICKETS, params: data, method: 'GET' })
            .then(resp => {
              resolve(resp)
              for (let i = 0; i < resp.data.length; i++) {
                const ticket = {
                  email: resp.data[i].email,
                  reply: resp.data[i]?.comments?.length,
                  ticketStatusLabel: this.getTicketStatus(resp.data.tickets[i].ticketStatus),
                  ticketStatus: resp.data.tickets[i].ticketStatus,
                  ticketId: resp.data[i].ticketId,
                  createdAtLabel: this.convertEpochToDateString(resp.data[i].createdAt),
                  createdAt: resp.data[i].createdAt,
                  title: resp.data[i].title,
                  projectId: resp.data[i].projectId,
                  projectName: resp.data[i].projectName,
                  ticketTypeLabel: this.getTicketType(resp.data[i].ticketType),
                  ticketType: resp.data[i].ticketType,
                  comments:  resp.data[i].comments,
                  body:  resp.data[i].body
                }
                this.tickets.push(ticket)
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
      if (this.token != "EOF" && this.ticketFilter == "") {
        this.readTickets()
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
    readTickets() {
      this.isLoading = true
      return new Promise((resolve, reject) => {
        let data = {
          email: this.$store.getters.email,
          pageSize: this.pageSize,
          token: this.token
        }
        axios({url: API_ENDPOINTS.ADMIN_READ_TICKETS, data: data, method: 'POST' })
            .then(resp => {
              resolve(resp)
              this.token = resp.data.token
              for (let i = 0; i < resp.data.tickets.length; i++) {
                const ticket = {
                  email: resp.data.tickets[i].email,
                  reply: resp.data.tickets[i]?.comments.length,
                  ticketStatusLabel: this.getTicketStatus(resp.data.tickets[i].ticketStatus),
                  ticketStatus: resp.data.tickets[i].ticketStatus,
                  ticketId: resp.data.tickets[i].ticketId,
                  createdAtLabel: this.convertEpochToDateString(resp.data.tickets[i].createdAt),
                  createdAt: resp.data.tickets[i].createdAt,
                  title: resp.data.tickets[i].title,
                  projectId: resp.data.tickets[i].projectId,
                  projectName: resp.data.tickets[i].projectName,
                  ticketTypeLabel: this.getTicketType(resp.data.tickets[i].ticketType),
                  ticketType: resp.data.tickets[i].ticketType,
                  comments: resp.data.tickets[i].comments,
                  body: resp.data.tickets[i].body
                }
                this.tickets.push(ticket)
              }
              this.isLoading = false
            })
            .catch(err => {
              reject(err)
              this.isLoading = false
            })
      })
    },
    getTicketStatus(ticketStatus) {
      switch(ticketStatus) {
        case 0:
          return "None"
        case 1:
          return "Open"
        case 2:
          return "Closed"
        case 3:
          return "Resolved"
        case 4:
          return "Deleted"
        default:
          return "None"
      }
    },
    getTicketType(ticketType) {
      switch(ticketType) {
        case 0:
          return "General"
        case 1:
          return "Bug"
        case 2:
          return "Payment"
        case 3:
          return "API"
        case 4:
          return "Feature"
        default:
          return "None"
      }
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
    if (!this.tickets || this.tickets.length == 0) {
      this.tickets = []
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

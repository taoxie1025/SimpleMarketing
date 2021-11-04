<template>
  <div>
    <v-expansion-panel-header>
      <v-row no-gutters>
        <v-col cols="8">
          <v-row>
            <v-col>
              <h3>{{ ticket.title }}</h3>
              <h6><i>{{getMetadata()}}</i></h6>
            </v-col>
          </v-row>
        </v-col>
        <v-spacer></v-spacer>
        <v-col cols="2" class="mr-3">
          <v-row>
            <v-col>
              <v-select
                  v-model="ticket.ticketStatus"
                  placeholder="Default"
                  :items="ticketStatus"
                  item-text="label"
                  item-value="value"
                  label="Status"
                  outlined
                  dense
                  @click.stop=""
                  @change="changeTicketStatus"
              ></v-select>
            </v-col>
          </v-row>
        </v-col>
        <v-col cols="1">
          <v-row>
            <v-col>
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn icon v-on:click.stop="openDialog" v-bind="attrs" v-on="on" class="pa-0 ma-0">
                    <v-icon size="35" color="grey">delete</v-icon>
                  </v-btn>
                </template>
                <span>Delete</span>
              </v-tooltip>
            </v-col>
            <v-col>
              <v-dialog v-model="dialog" max-width="250">
                <confirm-dialog v-bind:title="`Are you sure?`" :body="`The action is not reversible.`" @yes="deleteTicket(ticket.ticketId)" @no="dialog=false"></confirm-dialog>
              </v-dialog>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-expansion-panel-header>
    <v-expansion-panel-content>
      <v-row class="text--secondary ml-1 mr-1 mt-1">
        <p v-html="ticket.body"></p>
      </v-row>
      <hr>
      <v-row>
        <comments v-model="ticket.comments" :comments="ticket.comments" id="comments"></comments>
      </v-row>
      <v-row class="mt-2">
        <v-flex xs12>
          <tiptap-vuetify v-model="comment.body" :extensions="extensions" placeholder="Please reply here..." min-height="150"/>
        </v-flex>
      </v-row>
      <v-row>
        <v-btn outlined small tile class="mt-2" @click="cancel">Cancel</v-btn>
        <v-btn small tile color="info" class="mt-2 ml-2" @click="appendComment" :disabled="!isCommentValid">Reply</v-btn>
      </v-row>
    </v-expansion-panel-content>
  </div>
</template>

<script>

import { TiptapVuetify, Heading, Bold, Italic, Strike, Underline, Code, Paragraph, BulletList, OrderedList, ListItem, Link, Blockquote, HorizontalRule, History, CodeBlock, HardBreak, Image } from 'tiptap-vuetify'
import FileSelector from './Gallery'
import ImageForm from './ImageForm'
import Comments from './Comments'
import {EventBus as bus} from "../event_bus.js"
import ConfirmDialog from "./ConfirmDialog"
import axios from "axios"
import {default as API_ENDPOINTS} from "@/api"

export default {
  name: "TicketBrief",
  components: {TiptapVuetify, Comments, ConfirmDialog},
  props: ['ticket'],
  data() {
    return {
      dialog: false,
      isUpdating: false,
      comment: {
        email: this.$store.getters.email,
        ticketId: this.ticket.ticketId,
        name: "",
        body: ""
      },
      ticketStatus: [
        {label: "OPEN", value: 1},
        {label: "CLOSED", value: 2},
        {label: "RESOLVED", value: 3},
      ],
      extensions: [
        History,
        Blockquote,
        Link,
        Underline,
        Strike,
        Italic,
        ListItem,
        BulletList,
        OrderedList,
        [Heading, {
          options: {
            levels: [1, 2, 3]
          }
        }],
        Bold,
        Code,
        CodeBlock,
        HorizontalRule,
        Paragraph,
        HardBreak,
        [Image, {
          options: {
            imageSources: [
              { component: FileSelector, name: 'Uploads' },
              { component: ImageForm, name: 'Add By Link' }
            ],
            imageSourcesOverride: true
          }
        }],
      ],
    }
  },
  computed: {
    isCommentValid() {
      if (!this.comment.body || this.comment.body == "" || this.comment.body == "<p></p>") {
        return false
      }
      return true
    }
  },
  methods: {
    openDialog() {
      this.dialog = true
    },
    deleteTicket(ticketId) {
      return new Promise((resolve, reject) => {
        this.isUpdating = true
        axios({url: API_ENDPOINTS.DELETE_TICKET(this.$store.getters.email, ticketId), method: 'DELETE'})
            .then(resp => {
              resolve(resp)
              this.isUpdating = false
              this.$emit('ticketDeleted', ticketId)
              this.dialog = false
            })
            .catch(err => {
              reject(err)
              this.isUpdating = false
              this.dialog = false
              this.showQuickMessage("Failed to delete the ticket, please try again later.")
            })
      })
    },
    changeTicketStatus() {
      return new Promise((resolve, reject) => {
        this.isUpdating = true
        let data = {
          email: this.$store.getters.email,
          ticket: this.ticket
        }
        axios({url: API_ENDPOINTS.UPDATE_TICKET(this.ticket.ticketId), data: data, method: 'PUT'})
            .then(resp => {
              resolve(resp)
              this.isUpdating = false
              this.$emit('ticketUpdated', resp.data)
            })
            .catch(err => {
              reject(err)
              this.isUpdating = false
              this.showQuickMessage("Failed to update the ticket, please try again later.")
            })
      })
    },
    getMetadata() {
      let res = ""
      if (this.ticket?.comments?.length && this.ticket?.comments?.length > 0) {
        res += this.ticket?.comments.length + (this.ticket?.comments.length  > 1 ? " replies" : " reply")
      }
      if (res && res != "" && this.ticket?.projectName && this.ticket?.projectName != "") {
        res = res + " | " + this.ticket?.projectName
      } else {
        res += this.ticket?.projectName
      }
      return res
    },
    cancel() {
      this.$emit('cancel')
    },
    appendComment() {
      this.$emit('appendComment', this.comment)
    },
    onCommentAppended() {
      this.comment.body = ""
    }
  },
  mounted() {
    bus.$on('commentAppended', this.onCommentAppended)
  },
  created() {

  }
}
</script>

<style scoped>
#comments {
  height: 250px;
}
</style>
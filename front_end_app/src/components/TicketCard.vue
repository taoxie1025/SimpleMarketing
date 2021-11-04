<template>
  <v-card outlined tile >
    <v-card-text>
      <v-layout>
        <v-row align="center">
          <v-col class="d-flex" cols="12" sm="6">
            <v-select
              v-model="ticket.ticketType"
              placeholder="Default"
              :items="ticketTypes"
              item-text="label"
              item-value="value"
              label="Issue"
              outlined
              dense
              >
            </v-select>
          </v-col>
          <v-col class="d-flex" cols="12" sm="6">
            <v-select
                v-model="ticket.projectId"
                placeholder="None"
                :items="projects"
                item-text="option"
                item-value="value"
                label="Project(optional)"
                outlined
                dense
            ></v-select>
          </v-col>
        </v-row>
      </v-layout>
      <v-text-field v-model="ticket.title" label="Title" outlined required dense></v-text-field>
      <tiptap-vuetify v-model="ticket.body" :extensions="extensions" placeholder="Please describe your issue as detailed as possible..." min-height="150"/>
    </v-card-text>
    <v-card-actions>
      <v-btn outlined tile @click="cancel">Cancel</v-btn>
      <v-btn tile color="info" @click="createTicket" :disabled="!isTicketValid">Create</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>

import { TiptapVuetify, Heading, Bold, Italic, Strike, Underline, Code, Paragraph, BulletList, OrderedList, ListItem, Link, Blockquote, HorizontalRule, History, CodeBlock, HardBreak, Image } from 'tiptap-vuetify'
import FileSelector from './Gallery'
import ImageForm from './ImageForm'

export default {
  name: "TicketCard",
  components: {TiptapVuetify},
  props: [
    'ticket', 'projects'
  ],
  data() {
    return {
      ticketTypes: [
        {label: "General", value: 0},
        {label: "Payment", value: 1},
        {label: "Bug Report", value: 2},
        {label: "API Support", value: 3},
        {label: "Feature Request", value: 4}
      ],
      ticketStatus: [
        {label: "OPEN", value: 0},
        {label: "CLOSE", value: 1},
        {label: "RESOLVED", value: 2},
        {label: "DELETED", value: 3},
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
    isTicketValid() {
      if (!this.ticket.title || this.ticket.title == "") {
        return false
      }
      return true
    }
  },
  methods: {
    cancel() {
      this.$emit('cancel')
    },
    createTicket() {
      this.$emit('createTicket', this.ticket.projectId)
    }
  },
  created() {
    if (!this.ticket.ticketType) {
      this.ticket.ticketType = this.ticketTypes[0].value
    }
  }
}
</script>

<style scoped>

</style>
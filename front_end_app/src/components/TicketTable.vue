<template>
  <v-container>
    <v-layout row wrap>
      <v-expansion-panels popout focusable id="panel" v-model="panel">
        <v-expansion-panel v-for="ticket in tickets" :key="ticket.ticketId">
          <ticket-brief
              :ticket="ticket"
              @appendComment="onAppendComment"
              @cancel="onCancel"
              @ticketUpdated="onTicketUpdated"
              @ticketDeleted="onTicketDeleted">
          </ticket-brief>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-layout>
  </v-container>
</template>

<script>

import TicketBrief from "@/components/TicketBrief";
import {EventBus as bus} from "@/event_bus";

export default {
  components: {TicketBrief},
  props: ['tickets'],
  data() {
    return {
      panel: [],
    }
  },
  watch: {
  },
  computed: {
  },
  methods: {
    onTicketDeleted(ticketId) {
      for (let i = 0; i < this.tickets.length; i++) {
        if (this.tickets[i].ticketId == ticketId) {
          this.tickets.splice(i, 1)
          this.collapseAll()
          break
        }
      }
    },
    onTicketUpdated(updatedTicket) {
      for (let i = 0; i < this.tickets.length; i++) {
        if (this.tickets[i].ticketId == updatedTicket.ticketId) {
          this.tickets.splice(i, 1, updatedTicket)
          break
        }
      }
    },
    onCreatingNewTicket() {
      this.collapseAll()
    },
    collapseAll() {
      this.panel = []
    },
    onCancel() {
      this.collapseAll()
    },
    onAppendComment(comment) {
      this.$emit('appendComment', comment)
    }
  },
  mounted() {
    bus.$on('creatingNewTicket', this.onCreatingNewTicket)
  }
}
</script>

<style scoped>

</style>
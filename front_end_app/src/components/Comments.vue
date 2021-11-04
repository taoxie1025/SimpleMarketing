<template>
  <v-layout>
    <v-flex xs12 class="chat-area">
      <div v-for="(comment, i) in comments" :key="i">
        <message :comment="comment" :me="comment.email"></message>
      </div>
    </v-flex>
  </v-layout>
</template>

<script>

import Message from './Message'

export default {
  name: "Comments",
  components: {Message},
  props: [
    'comments'
  ],
  data() {
    return {

    }
  },
  watch: {
    comments: function () {
      this.scrollToBottom()
    }
  },
  computed: {
  },
  methods: {
    scrollToBottom() {
      let container = this.$el.querySelector("#commentWindow")
      if (container) {
        container.scrollTop = container.scrollHeight
      }
    }
  },
  created() {
  }
}
</script>

<style scoped>
.chat-area {
  /*   border: 1px solid #ccc; */
  max-height: 200px;
  background: white;
  padding: 1em;
  overflow-y: auto;
  margin: 0 auto 2em auto;
}
</style>
<template>
  <v-layout>
    <v-flex xs12>
      <section ref="chatArea">
        <p v-html="comment.body" class="message" :class="{ 'message-out': isMe, 'message-in ml-3': !isMe }"></p>
      </section>
    </v-flex>
  </v-layout>
</template>

<script>
// chat box style example
// https://www.csscodelab.com/vue-js-chat-box-code-example/
export default {
  name: "Message",
  components: {},
  props: [
    'comment', 'me'
  ],
  data() {
    return {

    }
  },
  computed: {
    isMe() {
      return this.me === this.$store.getters.email
    }
  },
  methods: {

  },
  created() {

  }
}
</script>

<style scoped>

body, html {
  font-weight: 100;
  background: #7b4397;  /* fallback for old browsers */
  background: -webkit-linear-gradient(to right, #dc2430, #7b4397);  /* Chrome 10-25, Safari 5.1-6 */
  background: linear-gradient(to right, #dc2430, #7b4397); /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */
}

.message {
  width: 45%;
  border-radius: 10px;
  padding: .3em;
  /*   margin-bottom: .3em; */
  font-size: .8em;
}
.message-out {
  background: #5bc0de;
  color: white;
  margin-left: 50%;
}
.message-in {
  background: #F1F0F0;
  color: black;
}
</style>
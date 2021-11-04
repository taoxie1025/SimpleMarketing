<template>
  <v-card flat class="mx-auto">
    <v-card-title class="headline grey lighten-2">
      <v-icon left color="black">mdi-email-send-outline</v-icon>
      <h3>Send a Test Email</h3>
    </v-card-title>
    <v-container>
      <br>
      <v-flex xs12 sm12>
        <v-text-field
            label="Destination Email Address"
            placeholder="Enter an email"
            outlined
            prepend-inner-icon="mdi-pencil-plus-outline"
            v-model="emailAddress"
            :rules="emailRules"
            block
            required
        >
        </v-text-field>
      </v-flex>
      <v-row class="mx-auto">
        <v-btn class="ma-1" tile outlined color="info" @click="close">Cancel</v-btn>
        <v-spacer></v-spacer>
        <v-btn class="ma-1" tile color="info" @click="send" :disabled="!isEmailValid(emailAddress)">Send</v-btn>
      </v-row>
    </v-container>
  </v-card>
</template>

<script>

export default {
  name: "EmailAddressCard",
  components: {},
  props: [
    'emailAddress'
  ],
  data() {
    return {
      emailRules: [
        v => !!v || "Required",
        v => /.+@.+\..+/.test(v) || "E-mail must be valid"
      ],
    }
  },
  computed: {
  },
  methods: {
    close() {
      this.$emit('close')
    },
    send() {
      this.$emit('send', this.emailAddress)
    },
    isEmailValid(email) {
      if (!(/.+@.+\..+/.test(email))) {
        return false
      }
      return true
    }
  },
  created() {
  }
}
</script>

<style scoped>

</style>
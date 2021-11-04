<template>
  <v-app id="app">
    <v-container fluid>
      <v-row justify="center" align="center">
        <v-col cols="12" lg="8">
          <h2>Password recovery</h2>
        </v-col>
        <v-col cols="12" lg="8">
          <v-card>
            <v-card-text>
              Use the form below to reset the password for your account
            </v-card-text>
            <v-form ref="changePasswordForm" v-model="isChangePasswordValid" lazy-validation>
              <v-text-field class="ml-10 mr-10" v-model="newPassword" :rules="[rules.required, rules.min]" :disabled="isUpdating" color="blue-grey lighten-2" type="password" label="New password" hint="At least 4 characters" counter outlined required dense></v-text-field>
              <v-text-field class="ml-10 mr-10" v-model="newPasswordConfirm" :rules="[rules.required, passwordMatch]" :disabled="isUpdating" color="blue-grey lighten-2" type="password" label="Reenter new password" hint="At least 4 characters" counter outlined required dense></v-text-field>
              <v-btn class="ml-10 mb-3" tile :disabled="!isChangePasswordValid || isUpdating" @click="updatePassword">Save</v-btn>
            </v-form>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-snackbar v-model="snackbar" timeout="3000" @click="snackbar=false">
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
      </v-row>
    </v-container>
  </v-app>
</template>

<script>

import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  components: {},
  data() {
    return {
      isChangePasswordValid: false,
      newPassword: "",
      newPasswordConfirm: "",
      isUpdating: false,
      rules: {
        required: value => !!value || "Required.",
        min: v => (v && v.length >= 4) || "Min 4 characters",
      },
      token: "",
      snackbar: false,
      quickMessage: ""
    }
  },
  watch: {
  },
  methods: {
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
    },
    updatePassword() {
      if (this.$refs.changePasswordForm.validate()) {
        return new Promise((resolve, reject) => {
          this.isUpdating = true
          const data = {token: this.token, newPassword: this.newPassword, newPasswordConfirm: this.newPasswordConfirm}
          axios({url: API_ENDPOINTS.UPDATE_PASSWORD, data: data, method: 'POST'})
              .then(resp => {
                resolve(resp)
                this.isUpdating = false
                this.showQuickMessage("Update succeeded, please sign back in after the page redirects...")
                const self = this
                setTimeout(function() {
                  self.$router.push('/')
                }, 3000)

              })
              .catch(err => {
                reject(err)
                this.isUpdating = false
                this.$refs.changePasswordForm.reset();
                this.showQuickMessage("Update password has failed, please try again.")
              })
        })
      }
    },
    isTokenValid() {
      return new Promise((resolve, reject) => {
        this.isUpdating = true
        const data = {token: this.token}
        axios({url: API_ENDPOINTS.IS_TOKEN_VALID, data: data, method: 'POST'})
            .then(resp => {
              resolve(resp)
              this.isUpdating = false
            })
            .catch(err => {
              reject(err)
            })
      })
    }
  },
  computed: {
    passwordMatch() {
      return () => this.newPassword === this.newPasswordConfirm || "Password must match";
    }
  },
  mounted() {
  },
  created() {
    this.token = this.$route.query.token
    this.isTokenValid()
        .then(() => {})
        .catch(() => {
          this.$router.push('/')
        })
  }
}
</script>


<style>

</style>

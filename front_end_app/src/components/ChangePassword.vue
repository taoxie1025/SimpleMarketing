<template>
  <v-container>
    <v-card>
      <v-card-text>
        Use the form below to change the password for your account
      </v-card-text>
      <v-form ref="changePasswordForm" v-model="isChangePasswordValid" lazy-validation>
        <v-text-field class="ml-10 mr-10" v-model="currentPassword" :rules="[rules.required]" :disabled="isUpdating" color="blue-grey lighten-2" type="password" label="Current password" outlined required dense></v-text-field>
        <v-text-field class="ml-10 mr-10" v-model="newPassword" :rules="[rules.required, rules.min]" :disabled="isUpdating" color="blue-grey lighten-2" type="password" label="New password" hint="At least 4 characters" counter outlined required dense></v-text-field>
        <v-text-field class="ml-10 mr-10" v-model="newPasswordConfirm" :rules="[rules.required, passwordMatch]" :disabled="isUpdating" color="blue-grey lighten-2" type="password" label="Reenter new password" hint="At least 4 characters" counter outlined required dense></v-text-field>
        <v-btn class="ml-10 mb-3" tile :disabled="!isChangePasswordValid || isUpdating" @click="changePassword">Save Changes</v-btn>
      </v-form>
    </v-card>
    <v-layout>
      <v-row>
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
      </v-row>
    </v-layout>
  </v-container>
</template>

<script>
import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  components: {},
  data() {
    return {
      isUpdating: false,
      updatePasswordRequest: {
        email: this.$store.getters.email,
      },
      currentPassword: "",
      newPassword: "",
      newPasswordConfirm: "",
      rules: {
        required: value => !!value || "Required.",
        min: v => (v && v.length >= 4) || "Min 4 characters",
      },
      isChangePasswordValid: false,
      snackbar: false,
      quickMessage: "",
    }
  },
  computed: {
    passwordMatch() {
      return () => this.newPassword === this.newPasswordConfirm || "Password must match";
    },
  },
  methods: {
    reset() {
      this.currentPassword = ""
      this.newPassword = ""
      this.newPasswordConfirm = "";
    },
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
    },
    changePassword() {
      if (this.$refs.changePasswordForm.validate()) {
        return new Promise((resolve, reject) => {
          this.isUpdating = true
          this.updatePasswordRequest.currentPassword = this.currentPassword
          this.updatePasswordRequest.newPassword = this.newPassword
          this.updatePasswordRequest.newPasswordConfirm = this.newPasswordConfirm
          axios({url: API_ENDPOINTS.CHANGE_PASSWORD, data: this.updatePasswordRequest, method: 'PUT'})
              .then(resp => {
                resolve(resp)
                this.isUpdating = false
                this.reset()
                this.$refs.changePasswordForm.reset();
                this.showQuickMessage("Your password is changed successfully.")
              })
              .catch(err => {
                reject(err)
                this.isUpdating = false
                this.showQuickMessage("Failed to change password.")
              })
        })
      }
    }
  }
}
</script>

<style scoped>
</style>
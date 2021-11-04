<template>
  <v-container>
    <v-card>
      <v-card-text>
        Use the form below to change the contact info for your account
      </v-card-text>
      <v-form ref="updateBasicInfoForm" v-model="isChangeBasicInfoValid" lazy-validation>
        <v-text-field class="ml-10 mr-10" v-model="updateBasicInfoRequest.email" :rules="[rules.required]" disabled color="blue-grey lighten-2" label="Email" outlined required dense></v-text-field>
        <v-text-field class="ml-10 mr-10" v-model="userInfo.firstName" :rules="[rules.required]" :disabled="isUpdating" color="blue-grey lighten-2" label="First name" outlined required dense></v-text-field>
        <v-text-field class="ml-10 mr-10" v-model="userInfo.lastName" :rules="[rules.required]" :disabled="isUpdating" color="blue-grey lighten-2" label="Last name" outlined required dense></v-text-field>
        <v-text-field class="ml-10 mr-10" v-model="userInfo.phoneNumber" :disabled="isUpdating" color="blue-grey lighten-2" label="Phone number" outlined dense></v-text-field>
        <v-text-field class="ml-10 mr-10" v-model="userInfo.address" :disabled="isUpdating" color="blue-grey lighten-2" label="Address" outlined dense></v-text-field>
        <v-btn class="ml-10 mb-3" tile :disabled="!isChangeBasicInfoValid || isUpdating" @click="updateUserBasicInfo">Save Changes</v-btn>
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
      accountInfo: {},
      userInfo: {
        firstName: "",
        lastName: "",
        phoneNumber: ""
      },
      isUpdating: false,
      updateBasicInfoRequest: {
        email: this.$store.getters.email,
        userInfo: {
          email: this.$store.getters.email
        },
      },
      rules: {
        required: value => !!value || "Required.",
      },
      isChangeBasicInfoValid: false,
      snackbar: false,
      quickMessage: "",
    }
  },
  computed: {
  },
  methods: {
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
    },
    updateUserBasicInfo() {
      if (this.$refs.updateBasicInfoForm.validate()) {
        return new Promise((resolve, reject) => {
          this.isUpdating = true
          this.updateBasicInfoRequest.userInfo.lastName = this.userInfo.lastName
          this.updateBasicInfoRequest.userInfo.firstName = this.userInfo.firstName
          this.updateBasicInfoRequest.userInfo.phoneNumber = this.userInfo.phoneNumber
          this.updateBasicInfoRequest.userInfo.address = this.userInfo.address
          axios({url: API_ENDPOINTS.UPDATE_USER_INFO(this.$store.getters.email), data: this.updateBasicInfoRequest, method: 'PUT'})
              .then(resp => {
                resolve(resp)
                this.isUpdating = false
                this.userInfo = resp.data
                this.showQuickMessage("Your contact info is changed successfully.")
              })
              .catch(err => {
                reject(err)
                this.isUpdating = false
                this.showQuickMessage("Failed to change contact info.")
              })
        })
      }
    },
    readAccountInfo() {
      return new Promise((resolve, reject) => {
        this.isUpdating = true
        const data = {requesterEmail: this.$store.getters.email}
        axios({url: API_ENDPOINTS.READ_ACCOUNT_INFO(this.$store.getters.email), params: data, method: 'GET'})
            .then(resp => {
              resolve(resp)
              this.accountInfo = resp.data
              this.userInfo.firstName = resp.data.firstName
              this.userInfo.lastName = this.accountInfo.lastName
              this.userInfo.phoneNumber = this.accountInfo.phoneNumber
              this.userInfo.address = this.accountInfo.address
              this.isUpdating = false
            })
            .catch(err => {
              reject(err)
              this.isUpdating = false
            })
      })
    }
  },
  created() {
    this.readAccountInfo()
  }
}
</script>

<style scoped>
</style>
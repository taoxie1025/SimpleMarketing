<template>
  <v-card outlined tile>
    <v-card-title>
      <h3>Edit User</h3>
      <v-spacer></v-spacer>
      <h6 v-html="userAccountInfoCopy.email"></h6>
    </v-card-title>
    <v-card-text>
      <v-layout>
        <v-row align="center">
          <v-col class="d-flex" cols="12" sm="6">
            <v-select
                v-model="userAccountInfoCopy.subscriptionPlan"
                placeholder="Default"
                :items="plans"
                item-text="label"
                item-value="value"
                label="Plan"
                outlined
                dense
            ></v-select>
          </v-col>
          <v-col class="d-flex" cols="12" sm="6">
            <v-select
                v-model="userAccountInfoCopy.paymentStatus"
                placeholder="Default"
                :items="paymentStatus"
                item-text="label"
                item-value="value"
                label="Payment Status"
                outlined
                dense
            ></v-select>
          </v-col>
          <v-col class="d-flex" cols="12" sm="6">
            <v-select
                v-model="userAccountInfoCopy.isBlock"
                placeholder="Default"
                :items="blockStatus"
                item-text="label"
                item-value="value"
                label="Block"
                outlined
                dense
            ></v-select>
          </v-col>
          <v-col class="d-flex" cols="12" sm="6">
            <v-select
                v-model="userAccountInfoCopy.userScope"
                placeholder="Default"
                :items="userScopes"
                item-text="label"
                item-value="value"
                label="Scope"
                outlined
                dense
            ></v-select>
          </v-col>
          <v-col class="d-flex" cols="12" sm="6">
            <v-text-field
                outlined
                v-model.number="userAccountInfoCopy.emailUsageInCycle"
                label="Email Usage"
                dense
            ></v-text-field>
          </v-col>
        </v-row>
      </v-layout>
    </v-card-text>
    <v-card-actions>
      <v-btn tile depressed @click="cancel">Close</v-btn>
      <v-spacer></v-spacer>
      <v-btn tile depressed color="info" :disabled="isUpdating" @click="updateUser">Update</v-btn>
    </v-card-actions>
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
  </v-card>
</template>

<script>


import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  components: {},
  props: ['userAccountInfo'],
  data() {
    return {
      userAccountInfoCopy: {},
      isUpdating: false,
      quickMessage: "",
      snackbar: false,
      plans: [
        {label: "Free", value: 0},
        {label: "Pro", value: 1},
        {label: "Ultra", value: 2},
      ],
      paymentStatus: [
        {label: "Completed", value: 0},
        {label: "In Progress", value: 1},
        {label: "Failed", value: 2},
      ],
      blockStatus: [
        {label: "Yes", value: true},
        {label: "No", value: false},
      ],
      userScopes: [
        {label: "User", value: 0},
        {label: "Admin", value: 1},
      ]
    }
  },
  watch: {
  },
  methods: {
    cancel() {
      this.$emit('cancel')
    },
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
    },
    updateUser() {
      this.isUpdating = true
      return new Promise((resolve, reject) => {
        let data = {
          email: this.$store.getters.email,
          editedUserInfo: {
            email: this.userAccountInfoCopy.email,
            isBlock: this.userAccountInfoCopy.isBlock,
            userScope: this.userAccountInfoCopy.userScope,
            subscriptionPlan: this.userAccountInfoCopy.subscriptionPlan,
            emailUsageInCycle: this.userAccountInfoCopy.emailUsageInCycle,
            stripeCustomerId: this.userAccountInfoCopy.stripeCustomerId,
            paymentStatus: this.userAccountInfoCopy.paymentStatus
          }
        }
        axios({url: API_ENDPOINTS.ADMIN_UPDATE_USER, data: data, method: 'POST'})
            .then(resp => {
              resolve(resp)
              this.$emit('userUpdated', resp.data)
            })
            .catch(err => {
              reject(err)
              this.showQuickMessage("Failed to update the user, please try again later.")
              this.userAccountInfoCopy = JSON.parse(JSON.stringify(this.userAccountInfo))
              this.isUpdating = false
            })
      })
    }
  },
  created: function () {
    this.userAccountInfoCopy = JSON.parse(JSON.stringify(this.userAccountInfo))
  }
}
</script>


<style>

</style>

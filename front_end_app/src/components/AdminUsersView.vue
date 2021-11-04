<template>
  <v-app id="app">
    <v-container fluid>
      <v-layout row wrap class="ma-1">
        <v-flex xs12>
          <v-text-field
              label="Email Address"
              placeholder="Search for user"
              outlined
              prepend-inner-icon="mdi-magnify"
              :append-icon="emailFilter != '' ? 'mdi-close' : ''"
              v-model="emailFilter"
              @click:append="clear"
          >
          </v-text-field>
        </v-flex>
        <v-flex xs12>
          <v-data-table
              @click:row="rowClick"
              v-model="selectedUserInfos"
              :headers="headers"
              :items="userInfos"
              item-key="email"
              :loading="isLoading"
              loading-text="Loading..."
              hide-default-footer
              :server-items-length="pageSize"
              disable-sort
              single-select
              class="elevation-0"
              v-infinite-scroll="loadMore"
              infinite-scroll-disabled="busy"
              infinite-scroll-distance="10"
              infinite-scroll-immediate-check="false"
          >
          </v-data-table>
        </v-flex>
      </v-layout>
      <v-layout>
        <v-dialog v-if="adminUserManagementDialog" v-model="adminUserManagementDialog" width="700px">
          <admin-user-management :userAccountInfo="selectedUserInfos[0]" @userUpdated="onUserUpdated" @cancel="onCancel"></admin-user-management>
        </v-dialog>
      </v-layout>
    </v-container>
  </v-app>
</template>

<script>

import axios from "axios";
import {default as API_ENDPOINTS} from "../api";
import AdminUserManagement from "@/components/AdminUserManagement";

export default {
  components: {AdminUserManagement},
  data() {
    return {
      headers: [
        { text: 'Email', align: 'start', value: 'email', class: "grey lighten-4"},
        { text: 'Stripe ID', align: 'start', value: 'stripeCustomerId', class: "grey lighten-4"},
        { text: 'Quota Usage', align: 'start', value: 'emailUsageInCycle', class: "grey lighten-4"},
        { text: 'Plan', align: 'start', value: 'subscriptionPlanLabel', class: "grey lighten-4"},
        { text: 'Payment Status', align: 'start', value: 'paymentStatusLabel', class: "grey lighten-4"},
        { text: 'Cycle Reset', align: 'start', value: 'lastClearCycleTimeLabel', class: "grey lighten-4"},
        { text: 'Blocked', align: 'start', value: 'isBlockLabel', class: "grey lighten-4"},
        { text: 'Scope', align: 'end', value: 'userScopeLabel', class: "grey lighten-4"}
      ],
      token: "",
      pageSize: 15,
      userInfos: [],
      isLoading: false,
      selectedUserInfos: [],
      emailFilter: "",
      adminUserManagementDialog: false,
    }
  },
  watch: {
    emailFilter() {
      if (this.emailFilter == "") {
        this.userInfos = []
        this.readUsers()
      } else {
        this.userInfos = []
        this.token = ""
        this.searchUsers()
      }
    }
  },
  methods: {
    onCancel() {
      this.adminUserManagementDialog = false
    },
    onUserUpdated(updatedUser) {
      for (let i = 0; i < this.userInfos.length; i++) {
        if (this.userInfos[i].email == updatedUser.email) {
          const updatedUserInfo = {
            email: updatedUser.email,
            stripeCustomerId: updatedUser.stripeCustomerId,
            emailUsageInCycle: updatedUser.emailUsageInCycle,
            subscriptionPlanLabel: this.getPlan(updatedUser.subscriptionPlan),
            subscriptionPlan: updatedUser.subscriptionPlan,
            paymentStatusLabel: this.getPaymentStatus(updatedUser.paymentStatus),
            paymentStatus: updatedUser.paymentStatus,
            lastClearCycleTimeLabel: this.convertEpochToDateString(updatedUser.lastClearCycleTime),
            lastClearCycleTime: updatedUser.lastClearCycleTime,
            isBlockLabel: this.getBlockLabel(updatedUser.isBlock),
            isBlock: updatedUser.isBlock,
            userScopeLabel: this.getUserScope(updatedUser.userScope),
            userScope: updatedUser.userScope
          }
          this.userInfos.splice(i, 1, updatedUserInfo)
          this.adminUserManagementDialog = false
          break
        }
      }
    },
    rowClick(item) {
      this.selectedUserInfos = []
      this.selectedUserInfos.push(item)
      this.adminUserManagementDialog = true
    },
    highlightClickedRow(item) {
      const tr = item.target.parentNode;
      tr.classList.add('highlight');
    },
    clear() {
      this.emailFilter = ""
    },
    searchUsers() {
      this.isLoading = true
      this.userInfos = []
      this.selectedUserInfos = []
      return new Promise((resolve, reject) => {
        let data = {emailFilter: this.emailFilter}
        axios({url: API_ENDPOINTS.ADMIN_SEARCH_USERS, params: data, method: 'GET' })
            .then(resp => {
              resolve(resp)
              for (let i = 0; i < resp.data.length; i++) {
                const userInfo = {
                  email: resp.data[i].email,
                  stripeCustomerId: resp.data[i].stripeCustomerId,
                  emailUsageInCycle: resp.data[i].emailUsageInCycle,
                  subscriptionPlanLabel: this.getPlan(resp.data[i].subscriptionPlan),
                  subscriptionPlan: resp.data[i].subscriptionPlan,
                  paymentStatusLabel: this.getPaymentStatus(resp.data[i].paymentStatus),
                  paymentStatus: resp.data[i].paymentStatus,
                  lastClearCycleTimeLabel: this.convertEpochToDateString(resp.data[i].lastClearCycleTime),
                  lastClearCycleTime: resp.data[i].lastClearCycleTime,
                  isBlockLabel: this.getBlockLabel(resp.data[i].isBlock),
                  isBlock: resp.data[i].isBlock,
                  userScopeLabel: this.getUserScope(resp.data[i].userScope),
                  userScope: resp.data[i].userScope
                }
                this.userInfos.push(userInfo)
              }
              this.isLoading = false
            })
            .catch(err => {
              reject(err)
              this.isLoading = false
            })
      })
    },
    loadMore() {
      if (this.token != "EOF" && this.emailFilter == "") {
        this.readUsers()
      }
    },
    validateLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        console.log("not signed in ");
        this.$router.push('/')
        return false
      }
      return true
    },
    readUsers() {
      this.isLoading = true
      return new Promise((resolve, reject) => {
        let data = {
          email: this.$store.getters.email,
          pageSize: this.pageSize,
          token: this.token
        }
        axios({url: API_ENDPOINTS.ADMIN_READ_USERS, data: data, method: 'POST' })
            .then(resp => {
              resolve(resp)
              this.token = resp.data.token
              for (let i = 0; i < resp.data.userAccountInfos.length; i++) {
                const userInfo = {
                  email: resp.data.userAccountInfos[i].email,
                  stripeCustomerId: resp.data.userAccountInfos[i].stripeCustomerId,
                  emailUsageInCycle: resp.data.userAccountInfos[i].emailUsageInCycle,
                  subscriptionPlanLabel: this.getPlan(resp.data.userAccountInfos[i].subscriptionPlan),
                  subscriptionPlan: resp.data.userAccountInfos[i].subscriptionPlan,
                  paymentStatusLabel: this.getPaymentStatus(resp.data.userAccountInfos[i].paymentStatus),
                  paymentStatus: resp.data.userAccountInfos[i].paymentStatus,
                  lastClearCycleTimeLabel: this.convertEpochToDateString(resp.data.userAccountInfos[i].lastClearCycleTime),
                  lastClearCycleTime: resp.data.userAccountInfos[i].lastClearCycleTime,
                  isBlockLabel: this.getBlockLabel(resp.data.userAccountInfos[i].isBlock),
                  isBlock: resp.data.userAccountInfos[i].isBlock,
                  userScopeLabel: this.getUserScope(resp.data.userAccountInfos[i].userScope),
                  userScope: resp.data.userAccountInfos[i].userScope
                }
                this.userInfos.push(userInfo)
              }
              this.isLoading = false
            })
            .catch(err => {
              reject(err)
              this.isLoading = false
            })
      })
    },
    convertEpochToDateString(timestamp) {
      if (timestamp && timestamp > 0) {
        const date = new Date(timestamp)
        return date.toLocaleDateString()
      }
      return "-"
    },
    getUserScope(scope) {
      switch(scope) {
        case 0:
          return "User"
        case 1:
          return "Admin"
        case 2:
          return "SuperAdmin"
        default:
          return "None"
      }
    },
    getPlan(plan) {
      switch (plan) {
        case 0:
          return "Free"
        case 1:
          return "Pro"
        case 2:
          return "Ultra"
        default:
          return "None"
      }
    },
    getPaymentStatus(paymentStatus) {
      switch (paymentStatus) {
        case 0:
          return "Completed"
        case 1:
          return "In Progress"
        case 2:
          return "Failed"
        default:
          return "None"
      }
    },
    getBlockLabel(isBlock) {
      return isBlock ? "Yes" : "No"
    }
  },
  created: function () {
    this.validateLogin()

    //bypass javascript not enabled issue
    if (!this.userInfos || this.userInfos.length == 0) {
      this.userInfos = []
      this.loadMore()
    }
  }
}
</script>


<style>
.v-text-field .v-icon {
  color: grey !important;
}
.v-data-table .v-icon {
  color: grey !important;
}
</style>

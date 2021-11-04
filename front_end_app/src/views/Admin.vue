<template>
  <v-app id="app">
    <v-container fluid>
      <v-tabs v-model="tabs" align-with-title grow light>
        <v-tab ripple>Users</v-tab>
        <v-tab ripple>Tickets</v-tab>
        <v-tab ripple>Subscribers</v-tab>
        <v-tab-item>
          <admin-users-view></admin-users-view>
        </v-tab-item>
        <v-tab-item>
          <v-card flat>
            <v-card-text>
              <admin-tickets-view></admin-tickets-view>
            </v-card-text>
          </v-card>
        </v-tab-item>
        <v-tab-item>
          <admin-subscribers-view></admin-subscribers-view>
        </v-tab-item>
      </v-tabs>
    </v-container>
  </v-app>
</template>

<script>

import AdminUsersView from "@/components/AdminUsersView";
import AdminTicketsView from "@/components/AdminTicketsView";
import AdminSubscribersView from "@/components/AdminSubscribersView";

export default {
  components: {AdminUsersView, AdminTicketsView, AdminSubscribersView},
  data() {
    return {
      items: ["Users", "Tickets"],
      accountInfo: {},
      tabs: []
    }
  },
  watch: {
  },
  methods: {
    validateLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        this.$router.push('/')
        return false
      }
      return true
    }
  },
  mounted() {
  },
  created: function () {
    this.validateLogin()
  }
}
</script>


<style>

</style>

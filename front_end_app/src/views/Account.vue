<template>
  <v-app id="app">
    <v-container fluid>
      <v-row justify="center" align="center">
        <v-col cols="12" lg="8">
          <h2>Change password</h2>
        </v-col>
        <v-col cols="12" lg="8">
          <change-password></change-password>
        </v-col>
        <v-col cols="12" lg="8" class="mt-3">
          <h2>Update contact info</h2>
        </v-col>
        <v-col cols="12" lg="8">
          <user-basic-info></user-basic-info>
        </v-col>
      </v-row>
    </v-container>
  </v-app>
</template>

<script>
import ChangePassword from "@/components/ChangePassword";
import UserBasicInfo from "@/components/UserBasicInfo";

export default {
  components: {UserBasicInfo, ChangePassword},
  data() {
    return {
      accountInfo: {}
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

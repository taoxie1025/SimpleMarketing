<template>
  <v-app id="app">
    <v-container fluid>
      <v-layout row wrap>
        <v-flex xs12 sm12 md6 lg4 v-for="(plan, index) in plans" v-bind:key="plan.value">
          <div @click="changeSelectedPlan(index)">
            <plan-card :plan="plan" :selectedPlan="selectedPlan" class="ma-3"></plan-card>
          </div>
        </v-flex>
      </v-layout>
      <v-layout>
        <v-row justify="center">
          <v-dialog persistent v-model="checkoutDialog" max-width="350">
            <v-card>
              <v-card-title class="headline grey lighten-2">
                <h3 v-if="selectedPlan>0">Change to {{plans[selectedPlan].name}} plan?</h3>
                <h3 v-else>Are you sure to unsubscribe?</h3>
              </v-card-title>
              <v-layout>
                <checkout class="mt-2" :plan="plans[selectedPlan]" :accountInfo="accountInfo" @closeDialog="checkoutDialog=false"></checkout>
              </v-layout>
              <v-container>
                <v-row class="mx-auto">
                  <v-btn class="ma-1" tile outlined color="info" @click="cancel">Cancel</v-btn>
                  <v-spacer></v-spacer>
                  <v-btn v-if="selectedPlan>0 && accountInfo.subscriptionId==''" class="ma-1" dark tile color="info" @click="checkout">Checkout</v-btn>
                  <v-btn v-else-if="selectedPlan>0 && accountInfo.subscriptionId!=''" class="ma-1" dark tile color="info" @click="checkout">Change</v-btn>
                  <v-btn v-else class="ma-1" tile outlined color="info" @click="unsubscribe">Yes</v-btn>
                </v-row>
              </v-container>
            </v-card>
          </v-dialog>
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
      </v-layout>
    </v-container>
  </v-app>
</template>
<script>

import PlanCard from "@/components/PlanCard";
import Checkout from "@/components/CheckoutComfirm";
import {EventBus as bus} from "../event_bus.js"
import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  components: {Checkout, PlanCard},
  data() {
    return {
      plans: [
        {name: "Free", price: "0", quota: "1000", value: 0,
          icon: require("@/assets/images/logos/free-plan.png"),
          feature: [
            {name: "Landing page", enabled: true},
            {name: "API & extension support", enabled: false},
            {name: "Developer support", enabled: false}
          ]
        },
        {name: "Pro", price: "20", quota: "100,000", value: 1, priceId: process.env.VUE_APP_PRO_PLAN_PRICE_ID,
          icon: require("@/assets/images/logos/pro-plan.png"),
          feature: [
            {name: "Landing page", enabled: true},
            {name: "API & extension support", enabled: true},
            {name: "Developer support", enabled: true},
          ]},
        {name: "Ultra", price: "60", quota: "300,000", value: 2, priceId: process.env.VUE_APP_ULTRA_PLAN_PRICE_ID,
          icon: require("@/assets/images/logos/ultra-plan.png"),
          feature: [
            {name: "Landing page", enabled: true},
            {name: "API & extension support", enabled: true},
            {name: "Developer support", enabled: true}
          ]}
      ],
      currentPlanInDb: -1,
      selectedPlan: -1,
      checkoutDialog: false,
      items: [
        {
          sku: 'prod_I7WQyekObdxoMN',
          quantity: 1
        }
      ],
      stripe: {},
      stripePublicKey: process.env.VUE_APP_STRIPE_PUBLIC_KEY,
      accountInfo: {},
      snackbar: false,
      quickMessage: "",
    }
  },
  watch: {
  },
  methods: {
    showQuickMessage(msg) {
      this.quickMessage = msg
      this.snackbar = true
    },
    unsubscribe() {
      this.cancelSubscription()
    },
    cancelSubscription() {
      return new Promise((resolve, reject) => {
        const data = {requesterEmail: this.$store.getters.email}
        axios({url: API_ENDPOINTS.CANCEL_SUBSCRIPTION(this.$store.getters.email, this.accountInfo.subscriptionId), params: data, method: 'DELETE'})
            .then(resp => {
              resolve(resp)
              this.readAccountInfo()
              this.checkoutDialog = false
            })
            .catch(err => {
              reject(err)
            })
      })
    },
    cancel() {
      this.selectedPlan = this.currentPlanInDb
      this.checkoutDialog = false
    },
    checkout() {
      bus.$emit('proceedToCheckout')
    },
    changeSelectedPlan(index) {
      this.selectedPlan = index
      if (this.selectedPlan == this.currentPlanInDb) {
        return
      }
      this.checkoutDialog = true
    },
    validateLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        this.$router.push('/')
        return false
      }
      return true
    },
    readAccountInfo() {
      return new Promise((resolve, reject) => {
        const data = {requesterEmail: this.$store.getters.email}
        axios({url: API_ENDPOINTS.READ_ACCOUNT_INFO(this.$store.getters.email), params: data, method: 'GET'})
            .then(resp => {
              resolve(resp)
              this.accountInfo = resp.data
              this.currentPlanInDb = this.accountInfo.subscriptionPlan
              this.selectedPlan = this.currentPlanInDb
            })
            .catch(err => {
              reject(err)
            })
      })
    }
  },
  created: function () {
    this.validateLogin()
    this.readAccountInfo()
    if (this.$route.query.status == "ok") {
      this.showQuickMessage("Payment succeeded.")
    } else if (this.$route.query.showStatus == "failed") {
      this.showQuickMessage("Payment failed, please try again.")
    }
  }
}
</script>


<style>

</style>

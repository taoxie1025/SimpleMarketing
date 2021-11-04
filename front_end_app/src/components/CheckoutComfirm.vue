<template>
  <div>
    <v-flex xs12>
      <v-row>
        <v-col class="ml-5">
        <img
            :src="plan.icon"
            width="100px"
        />
        </v-col>
        <v-col class="mt-5 ml-16">
          <h3>{{plan.name}}</h3>
          <h5 v-if="shownPrice==plan.price">${{plan.price}} / month</h5>
          <div v-else>
            <h5><strike>${{plan.price}} / month</strike></h5>
            <h5 style="color:red;">${{shownPrice}} / month</h5>
          </div>
        </v-col>
      </v-row>
    </v-flex>
    <v-flex xs12>
      <v-row>
        <v-col class="ml-5">
          <v-text-field ref="couponRef" @click:append="clearCoupon" :append-icon="getIcon" outlined dense label="Coupon" v-model="coupon" style="width:195px"></v-text-field>
        </v-col>
        <v-col>
          <v-btn outlined tile @click="applyCoupon" :disabled="!coupon">Apply</v-btn>
        </v-col>
      </v-row>
    </v-flex>
    <v-snackbar v-model="snackbar" timeout="4000" @click="snackbar=false">
      {{msg}}
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
  </div>
</template>
<script>

import {loadStripe} from '@stripe/stripe-js';
import {EventBus as bus} from "../event_bus.js"
import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

export default {
  components: {},
  props: ['plan', 'accountInfo'],
  data() {
    return {
      stripe: {},
      elements: {},
      card: {},
      stripePublicKey: process.env.VUE_APP_STRIPE_PUBLIC_KEY,
      snackbar: false,
      msg: "",
      coupon: "",
      isCouponValid: false,
      shownPrice: ""
    }
  },
  watch: {
    coupon: function () {
      this.isCouponValid = false
      this.shownPrice = this.plan.price
    },
    plan: function() {
      this.shownPrice = this.plan.price
      this.isCouponValid = false
      this.coupon = ""
    }
  },
  computed: {
    getIcon() {
      if (!this.coupon) {
        return "mdi-ticket-percent-outline"
      } else if (this.isCouponValid) {
        return "mdi-check"
      }
      return "mdi-close"
    }
  },
  methods: {
    clearCoupon() {
      if (!this.isCouponValid && this.coupon) {
        this.coupon = ""
        this.isCouponValid = false
        this.shownPrice = this.plan.price
      }
    },
    applyCoupon() {
      this.readCoupon()
          .then((resp) => {
            this.isCouponValid = true
            this.shownPrice = resp.data.percent_off / 100 * this.plan.price
          })
          .catch(() => {
            this.isCouponValid = false
            this.coupon = ""
            this.$refs.couponRef.focus()
            this.showErrorMessage("Coupon is not valid")
          })
    },
    readCoupon() {
      return axios({
        url: API_ENDPOINTS.READ_COUPON(this.coupon),
        method: 'GET'
      })
    },
    updateSubscription() {
      const data = {
        priceId: this.plan.priceId,
      }
      return axios({
        url: API_ENDPOINTS.UPDATE_SUBSCRIPTION(this.accountInfo.email, this.accountInfo.subscriptionId),
        data: data,
        method: 'PUT'
      })
    },
    createSession() {
      const data = {
        coupon: this.coupon,
        priceId: this.plan.priceId,
        stripeCustomerId: this.accountInfo.stripeCustomerId,
        plan: this.plan.name,
        successUrl: process.env.VUE_APP_WEB_HOST +"/plan?status=ok",
        cancelUrl: process.env.VUE_APP_WEB_HOST + "/plan?status=failed",
        email: this.$store.getters.email ? this.$store.getters.email : this.accountInfo.email
      }
      return axios({url: API_ENDPOINTS.CHECKOUT_CREATE_SESSION, data: data, method: 'POST' })
    },
    checkout() {
      if (this.accountInfo?.subscriptionPriceId == "") {
        // if the user has no subscriptionPriceId, it is a brand new subscription.
        this.createSession().
        then((resp) => {
          this.stripe.redirectToCheckout({ sessionId: resp.data.id })
        }).
        catch(() => {
          this.showErrorMessage("Failed to create checkout session, please try again.")
        })
      } else {
        // if the user has a subscriptionPriceId, it means that the user has already subscribed.
        // so here we just need to update the priceId in the user current subscription.
        // price will be adjusted on a proration basis by Stripe. https://stripe.com/docs/api/subscriptions/update
        this.updateSubscription().
            then(() => {
              this.$emit("closeDialog")
        }).catch(() => {
          this.showErrorMessage("failed to change plan, please try again later")
        })
      }
    },
    showErrorMessage(msg) {
      this.msg = msg
      this.snackbar = true
    }
  },
  mounted() {
    bus.$on('proceedToCheckout', this.checkout)
  },
  created: async function () {
    this.stripe = await loadStripe(this.stripePublicKey)
    this.shownPrice = this.plan.price
  }
}
</script>


<style>

</style>

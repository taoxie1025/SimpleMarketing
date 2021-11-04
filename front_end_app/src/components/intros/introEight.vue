<template>
  <!-- intro-section -->
  <section id="intro-wrap" class="intro-eight-Wrap text-white  text-center">
    <b-container >
      <b-row>
        <b-col sm="12" >
          <div class="intro-eight-subtitle">
            <h3 class="font-weight-bold mb-3 text-white">
              Simple Marketing
            </h3>
          </div>
          <h1
            class="intro-eight-Title font-weight-bold text-42 t-shadow mb-3 text-white"
          >
            Be Productive While Marketing Your Business
          </h1>
          <div class=" intro-eight-Description ">
            <p class="text-18 mb-4">
              Build Relationships That Influence Your Audience Using Our Intuitive Email Campaign Tool.
              You Only Focus on Your Contents and We Will Do The Heavy Lifting for You.
            </p>
            <p v-if="!isLogin">
              Get Started Today - No Credit Card Needed!
            </p>
          </div>

          <div v-if="!isLogin" class="intro-eight-Buttons mb-5">
            <a
              name=""
              id=""
              @click="openSignInDialog()"
              class="btn half-button btn-outline-white btn-lg pl-5 pr-5 pb-2  mr-2 mb-4 text-uppercase"
              role="button"
            >
              Sign in
            </a>
            <a
              id=""
              @click="openSignUpDialog()"
              class="btn half-button btn-warning btn-lg pl-5 pr-5 pb-2 mb-4 text-uppercase"
              role="button"
            >
              Sign up
            </a>
          </div>
          <div v-else class="intro-eight-Buttons mb-5">
            <a
                id=""
                @click="goToDashboard()"
                class="btn half-button btn-warning btn-lg pl-5 pr-5 pb-2 mb-4 text-uppercase"
                role="button"
            >
              Dashboard
            </a>
          </div>
        </b-col>
      </b-row>
    </b-container>
    <div class="overlay"></div>
    <v-container>
      <v-dialog v-model="authDialog" max-width="700px">
        <v-card color="blue-grey darken-1" dark>
          <auth :tab="selectedTabIndex"></auth>
        </v-card>
      </v-dialog>
    </v-container>
  </section>
  <!-- end::intro-section -->
</template>

<script>
import Auth from "../Auth.vue"
export default {
  name: 'Home',
  components: {
    Auth
  },
  data () {
    return {
      authDialog: false,
      selectedTabIndex: 0,
    }
  },
  methods: {
    openSignInDialog() {
      this.selectedTabIndex = 0
      this.authDialog = true
    },
    openSignUpDialog() {
      this.selectedTabIndex = 1
      this.authDialog = true
    },
    goToDashboard() {
      this.$router.push('/dashboard')
    }
  },
  watch: {
  },
  computed: {
    isLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        return false
      }
      return true
    },
  },
  created: function () {

  }
}
</script>
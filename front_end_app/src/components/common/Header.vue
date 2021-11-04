<template>
  <div class="main-header header-fixed-default" id="home-header">
    <b-navbar
      class="navbar container w-100 navbar-expand-lg navbar-transparent bg-transparent"
      toggleable="lg"
      type=""
      variant=""
    >
      <b-navbar-brand href="#">
        <div class="logo">
          <img src="@/assets/images/landing/lOGO_d/logo2.png" alt="" />
        </div>
      </b-navbar-brand>

      <b-navbar-toggle
        target="nav-collapse"
        class="eva eva-menu-outline text-18 text-white"
      >
      </b-navbar-toggle>

      <!-- <div class="menu-toggle navbar-toggler" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                    <div></div>
                    <div></div>
                    <div></div>
                </div> -->

      <b-collapse id="nav-collapse" is-nav>
        <div style="margin: auto"></div>
        <!-- Right aligned nav items -->
        <b-navbar-nav class="header-part-right">
          <!-- <b-nav-item href="#intro-wrap" v-smooth-scroll class="m-2"
            >Home <span class="sr-only">(current)</span></b-nav-item
          > -->

          <li class="nav-item">
            <a href="#intro-wrap" v-smooth-scroll class="m-2"
              >Home <span class="sr-only">(current)</span></a
            >
          </li>
          <li class="nav-item">
            <a href="#features-wrap" class="m-2" v-smooth-scroll>Features</a>
          </li>
          <li class="nav-item">
            <a href="#services-wrap" class="m-2" v-smooth-scroll>Services</a>
          </li>
          <li class="nav-item">
            <a href="#extra-feature-wrap" class="m-2" v-smooth-scroll>Examples</a>
          </li>
          <li class="nav-item">
            <a href="#pricing-wrap" class="m-2" v-smooth-scroll>Pricing</a>
          </li>
          <li class="nav-item b-nav-dropdown dropdown">
            <a href="#teams-wrap" target="_self" class="m-2" v-smooth-scroll>About</a>
          </li>
          <li class="nav-item">
            <a href="#contacts-wrap" class="m-2" v-smooth-scroll>Contact Us</a>
          </li>
          <li v-if="isLogin" class="nav-item">
            <a @click="goToDashboard()" class="m-2 btn half-button btn-warning" v-smooth-scroll>Dashboard</a>
          </li>
          <li v-else class="nav-item">
            <a @click="signIn()" class="m-2 btn half-button btn-warning" v-smooth-scroll>Sign In</a>
          </li>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <v-container>
      <v-dialog v-model="authDialog" max-width="700px">
        <v-card color="blue-grey darken-1" dark>
          <auth :tab="0"></auth>
        </v-card>
      </v-dialog>
    </v-container>
  </div>
</template>
<script>
import Auth from "../Auth.vue"

export default {
  components: {
    Auth
  },
  data() {
    return {
      link: {
        hash: ["#testimonials-wrap"]
      },
      authDialog: false
    };
  },
  methods: {
    handleScroll() {
      var scroll = window.pageYOffset;
      if (scroll >= 80) {
        document.querySelector(".main-header").classList.add("header-fixed");
      } else {
        document.querySelector(".main-header").classList.remove("header-fixed");
      }
    },
    goToDashboard() {
      this.$router.push('/dashboard')
    },
    signIn() {
      this.authDialog = true
    }
  },
  computed: {
    isLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        return false
      }
      return true
    },
  },
  created() {
    window.addEventListener("scroll", this.handleScroll);
  },
  destroyed() {
    window.removeEventListener("scroll", this.handleScroll);
  }
};
</script>

<style>

</style>

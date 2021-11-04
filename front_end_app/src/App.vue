<template>
  <v-app class="grey lighten-4">
    <v-main class="ma-0" v-if="isLandingPage || isRecoveryPage || isUnsubscribePage">
      <div class="landing_wrap">
        <router-view></router-view>
      </div>
    </v-main>
    <v-main class="mb-4" v-else-if="showNavBar">
      <Navbar></Navbar>
      <router-view></router-view>
    </v-main>
    <v-main class="ma-0" v-else>
      <div class="landing_wrap">
        <dx-header></dx-header>
        <router-view></router-view>
        <dx-contact></dx-contact>
        <dxFooter></dxFooter>
      </div>
    </v-main>
  </v-app>
</template>

<script>
  import Navbar from "./components/Navbar";
  import Header from "@/components/common/Header.vue";
  import Footer from "@/components/common/footer";
  import Contact from "@/components/common/contact";
  export default {
    name: 'App',
    components: {
      Navbar,
      dxHeader: Header,
      dxFooter: Footer,
      dxContact: Contact
    },
    data () {
      return {
      }
    },
    methods: {
      handleStorageChangeEvent (event) {
        if(event.key === 'logged_in') {
          if (!event.value) {
            this.$router.push('/')
          } else {
            location.reload()
          }
        }
      }
    },
    computed: {
      showNavBar() {
        return this.$route.path != '/';
      },
      isLandingPage() {
        return this.$route.path.includes('/p/')
      },
      isRecoveryPage() {
        return this.$route.path.includes('/recovery')
      },
      isUnsubscribePage() {
        return this.$route.path.includes('/unsubscribe')
      }
    },
    mounted() {
      window.addEventListener('storage', this.handleStorageChangeEvent, false)
    },
    created: function () {

    }
  }
</script>
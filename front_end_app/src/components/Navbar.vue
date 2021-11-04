<template>
    <nav>
        <v-app-bar text app>
            <v-app-bar-nav-icon class="grey--text" @click="drawer = !drawer"></v-app-bar-nav-icon>
            <v-toolbar-title v-if="!drawer" class="text-uppercase grey--text">
                <span>
                    <img width="200px" src="@/assets/images/landing/lOGO_d/logo3.png">
                </span>
            </v-toolbar-title>
            <v-spacer></v-spacer>
            <v-btn v-if="isLogin" text color="grey" @click="signout">
                <span>Sign out</span>
                <v-icon right>exit_to_app</v-icon>
            </v-btn>
        </v-app-bar>

        <!-- Drawer -->
        <v-navigation-drawer v-model="drawer" app class="primary">
            <v-layout column align-center>
                <v-flex class="mt-5 text-xs-center">
                  <v-toolbar-title class="text-uppercase white--text">
                    <span>
                      <img src="@/assets/images/landing/lOGO_d/logo7.png">
                    </span>
                  </v-toolbar-title>
                </v-flex>
                <v-col>
                  <v-divider class="mt-0 mb-20" color="grey"></v-divider>
                </v-col>
            </v-layout>
            <v-list>
                <v-list-item v-for="link in links" :key="link.text" router :to="link.route">
                    <v-list-item-action v-if="!link.hidden||isAdmin">
                        <v-icon class="white--text">{{ link.icon }}</v-icon>
                    </v-list-item-action>
                    <v-list-item-content v-if="!link.hidden||isAdmin">
                        <v-list-item-title class="white--text">{{ link.text }}</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>
    </nav>
</template>


<script>
    export default {
        components: {},
        data() {
            return {
                name: '',
                photoURL: '',
                drawer: true,
                links: [
                    { icon: 'dashboard', text: 'Dashboard', route: '/dashboard'},
                    { icon: 'mdi-alpha-p-box-outline', text: 'Plan', route: '/plan'},
                    { icon: 'mdi-face-agent', text: 'Support', route: '/support'},
                    { icon: 'mdi-account-cog', text: 'Account', route: '/account'},
                    { icon: 'mdi-monitor-lock', text: 'Admin', route: '/admin', hidden: true},
                ],
                snackbar: false,
            }
        },
        methods: {
          signout: function() {
            this.$store.dispatch('logout')
                .then(() => {
                  this.$router.push('/')
                })
          }
        },
        computed: {
          isLogin() {
            if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
              return false
            }
            return true
          },
          isAdmin() {
            if (this.$store.getters.userScope > 0) {
              return true
            }
            return false
          }
        },
        created() {
        }

    }
</script>
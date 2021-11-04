<template>

  <div id="app">
        <div>
          <v-tabs v-model="tab" show-arrows background-color="deep-purple accent-3" icons-and-text dark grow>
            <v-tabs-slider color="purple darken-4"></v-tabs-slider>
            <v-tab v-for="(tab,i) in tabs" :key="tab+i">
              <v-icon large>{{ tab.icon }}</v-icon>
              <div class="caption py-1">{{ tab.name }}</div>
            </v-tab>
            <v-tab-item>
              <v-card class="px-4">
                <v-card-text>
                  <v-form ref="loginForm" v-model="valid" lazy-validation>
                    <v-row>
                      <v-col cols="12">
                        <v-text-field v-model="loginEmail" :rules="loginEmailRules" label="E-mail" required></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-text-field v-if="!isPasswordRecover" v-model="loginPassword" :append-icon="show1?'eye':'eye-off'" :rules="[rules.required, rules.min]" :type="show1 ? 'text' : 'password'" name="input-10-1" label="Password" hint="At least 4 characters" counter @click:append="show1 = !show1"></v-text-field>
                      </v-col>
                      <v-col class="d-flex" cols="12" sm="3" xsm="12" align-end>
                        <v-btn x-small block text  color="blue" @click="showForgotPassword" v-html="!isPasswordRecover ? `Forgot Password?` : `Login`"></v-btn>
                      </v-col>
                      <v-spacer></v-spacer>
                      <v-col class="d-flex" cols="12" sm="3" xsm="12" align-end>
                        <v-btn v-if="!isPasswordRecover" :disabled="!valid" x-large block color="success" @click="validateLogin"> Login </v-btn>
                        <v-btn v-else x-large block :disabled="isUpdating || !isEmail(loginEmail)" color="success" @click="dialog=true"> Continue </v-btn>
                      </v-col>
                    </v-row>
                  </v-form>
                </v-card-text>
              </v-card>
            </v-tab-item>
            <v-tab-item>
              <v-card class="px-4">
                <v-card-text>
                  <v-form ref="registerForm" v-model="valid" lazy-validation>
                    <v-row>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field v-model="firstName" :rules="[rules.required]" label="First Name" maxlength="20" required></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="6">
                        <v-text-field v-model="lastName" :rules="[rules.required]" label="Last Name" maxlength="20" required></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-text-field v-model="email" :rules="emailRules" label="E-mail" required></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-text-field v-model="password" :append-icon="show1 ? 'mdi-eye' : 'mdi-eye-off'" :rules="[rules.required, rules.min]" :type="show1 ? 'text' : 'password'" name="input-10-1" label="Password" hint="At least 4 characters" counter @click:append="show1 = !show1"></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-text-field block v-model="verify" :append-icon="show1 ? 'mdi-eye' : 'mdi-eye-off'" :rules="[rules.required, passwordMatch]" :type="show1 ? 'text' : 'password'" name="input-10-1" label="Confirm Password" counter @click:append="show1 = !show1"></v-text-field>
                      </v-col>
                      <v-spacer></v-spacer>
                      <v-col class="d-flex ml-auto" cols="12" sm="3" xsm="12">
                        <v-btn x-large block :disabled="!valid" color="success" @click="validateSignUp">Register</v-btn>
                      </v-col>
                    </v-row>
                  </v-form>
                </v-card-text>
              </v-card>
            </v-tab-item>
          </v-tabs>
          <v-row justify="center">
            <v-dialog v-model="dialog" max-width="350">
              <confirm-dialog v-bind:title="`Password Recovery`" :body="`If you have an account with us, a link will be sent to your email to reset your password. The link will be expired in 1 hour.`" @yes="sendResetLink" @no="dialog=false"></confirm-dialog>
            </v-dialog>
          </v-row>
        </div>
  </div>

</template>

<script>
import ConfirmDialog from "./ConfirmDialog";
import axios from "axios";
import {default as API_ENDPOINTS} from "@/api";

  export default {
    components: {ConfirmDialog},
    props: [
      'tab'
    ],
    data() {
      return {
        isUpdating: false,
        isPasswordRecover: false,
        dialog: false,
        tabs: [
          {name:"Login", icon:"mdi-account"},
          {name:"Register", icon:"mdi-account-outline"}
        ],
        valid: true,

        firstName: "",
        lastName: "",
        email: "",
        password: "",
        verify: "",
        loginPassword: "",
        loginEmail: "",
        loginEmailRules: [
          v => !!v || "Required",
          v => /.+@.+\..+/.test(v) || "E-mail must be valid"
        ],
        emailRules: [
          v => !!v || "Required",
          v => /.+@.+\..+/.test(v) || "E-mail must be valid"
        ],

        show1: false,
        rules: {
          required: value => !!value || "Required.",
          min: v => (v && v.length >= 4) || "Min 4 characters"

        }
      }
    },

    methods: {
      isEmail(email) {
        const re = /\S+@\S+\.\S+/;
        return re.test(email);
      },
      sendResetLink() {
        return new Promise((resolve, reject) => {
          this.isUpdating = true
          const data = {email: this.loginEmail}
          axios({url: API_ENDPOINTS.RESET_PASSWORD, data: data, method: 'POST'})
              .then(resp => {
                resolve(resp)
                this.isUpdating = false
                this.dialog = false
              })
              .catch(err => {
                reject(err)
                this.isUpdating = false
                this.dialog = false
              })
        })
      },
      showForgotPassword() {
        this.isPasswordRecover = !this.isPasswordRecover
      },
      validateLogin() {
        if (this.$refs.loginForm.validate()) {
          this.signin();
        }
      },
      validateSignUp() {
        if (this.$refs.registerForm.validate()) {
          this.signup();
        }
      },
      reset() {
        //this.$refs.loginForm.reset();
        this.loginPassword = "";
        this.password = "";
      },
      resetValidation() {
        this.$refs.form.resetValidation();
      },
      signup() {
        let data = {
          firstName: this.firstName,
          lastName: this.lastName,
          email: this.email,
          password: this.password,
        }
        this.$store.dispatch('register', data)
          .then(() => {
            this.loginEmail = this.email
            this.loginPassword = this.password
            this.signin()
          })
          .catch(err => {
            if (err.length > 0) {
              this.reset()
            }
          })
      },

      signin() {
        let data = {
          email: this.loginEmail,
          password: this.loginPassword
        }
        this.$store.dispatch('login', data)
                .then(() => {
                  this.$router.push('/dashboard')
                })
                .catch(err => {
                  if (err.length > 0) {
                    this.reset()
                  }
                  this.reset()
                })

      },

    },
    computed: {
      passwordMatch() {
        return () => this.password === this.verify || "Password must match";
      }
    },
  }
</script>


<style>

  @media only screen and (max-width: 768px) {
    .v-content {
      margin: 0;
    }
  }

</style>

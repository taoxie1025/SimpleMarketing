<template>
  <!-- start::product-details-right-image -->
  <section id="pricing-wrap" class="pricing-four p-t-b-80">
    <b-container>
      <b-row>
        <b-col md="12">
          <div class="section-header col-md-12 mb-5">
            <h2 class="font-weight-bold">Choose Your Plan</h2>
            <p>
              Simple Marketing offers 7-day money back guarantee with no question ask.
            </p>
          </div>
        </b-col>
      </b-row>
      <b-row>
        <b-col md="6" lg="4" class="mb-lg-0  mb-4">
          <div class="pricing-four__card">
            <b-card>
              <div class="">
                <div class="pricing-four__contents">
                  <div class="pricing-four__titles text-center">
                    <h2 class="font-weight-bold mt-4 mb-4">Free</h2>
                  </div>
                  <div class="pricing-four__values  text-center mt-4 mb-4">
                    <p class="text-25">
                      $ <span class="text-50 font-weight-bold">0</span> / month
                    </p>
                  </div>

                  <ul class="list-group list-group-flush mt-5">
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                        class="eva eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span class="font-weight-800">1,000</span> emails per month
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                          class="eva eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span>Landing page</span>
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i class=" mr-4"></i>
                      <span class="text-mute"
                        ><strike>API & extension support</strike></span
                      >
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i class=" mr-4"></i>
                      <span class="text-mute"
                        ><strike>Developer support</strike></span
                      >
                    </li>
                  </ul>

                  <div v-if="!isLogin" class="text-center">
                    <b-button
                      variant="outline-primary "
                      class="mt-5 mb-4 pr-5 pl-5 btn-rounded"
                      @click="openSignUpDialog()"
                      >Sign Up</b-button
                    >
                  </div>
                </div>
              </div>
            </b-card>
          </div>
        </b-col>
        <b-col md="6" lg="4" class=" mb-lg-0  mb-4">
          <div class="pricing-four__card">
            <b-card class="overlay">
              <div class="">
                <div class="pricing-four__contents">
                  <div class="pricing-four__titles--white text-center">
                    <h2 class="font-weight-bold mt-4 mb-4 text-white">Pro</h2>
                  </div>
                  <div
                    class="pricing-four__values  text-center mt-4 mb-4 text-white"
                  >
                    <p class="text-25">
                      $ <span class="text-50 font-weight-bold">20</span> / month
                    </p>
                  </div>

                  <ul class="list-group list-group-flush mt-5 text-white">
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                        class="eva eva--white eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span class="font-weight-800">100,000</span> emails per month
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                          class="eva eva--white eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span>Landing page</span>
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                        class="eva eva--white eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span class="">API & extension support</span>
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                        class="eva eva--white eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span class="">Developer support</span>
                    </li>
                  </ul>

                  <div v-if="!isLogin" class="text-center">
                    <b-button
                      variant="outline-white"
                      class="btn-rounded mt-5 mb-4 pr-5 pl-5"
                      @click="openSignUpDialog()"
                      >Sign Up</b-button
                    >
                  </div>
                </div>
              </div>
            </b-card>
          </div>
        </b-col>
        <b-col md="6" lg="4" class="mb-lg-0">
          <div class="pricing-four__card">
            <b-card>
              <div class="">
                <div class="pricing-four__contents">
                  <div class="pricing-four__titles text-center">
                    <h2 class="font-weight-bold mt-4 mb-4">Ultra</h2>
                  </div>
                  <div class="pricing-four__values  text-center mt-4 mb-4">
                    <p class="text-25">
                      $ <span class="text-50 font-weight-bold">60</span> / month
                    </p>
                  </div>

                  <ul class="list-group list-group-flush mt-5">
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                        class="eva eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span class="font-weight-800">300,000</span> emails per month
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                          class="eva eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span>Landing page</span>
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                        class="eva eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span class="">API & extension support</span>
                    </li>
                    <li class="list-group-item bg-transparent border-0 pb-0">
                      <i
                        class="eva eva-checkmark-outline align-middle mr-2"
                      ></i>
                      <span class="">Developer support</span>
                    </li>
                  </ul>

                  <div v-if="!isLogin" class="text-center">
                    <b-button
                      variant="outline-primary"
                      class="btn-rounded mt-5 mb-4 pr-5 pl-5"
                      @click="openSignUpDialog()"
                      >Sign Up</b-button
                    >
                  </div>
                </div>
              </div>
            </b-card>
          </div>
        </b-col>
      </b-row>
    </b-container>
    <v-container>
      <v-dialog v-model="authDialog" max-width="700px">
        <v-card color="blue-grey darken-1" dark>
          <auth :tab="1"></auth>
        </v-card>
      </v-dialog>
    </v-container>
  </section>
  <!-- end::product-details-right-image -->
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
      authDialog: false
    }
  },
  methods: {
    openSignUpDialog() {
      this.authDialog = true
    },
  },
  computed: {
    isLogin() {
      if (!this.$store.getters.isLoggedIn || !this.$store.getters.email) {
        return false
      }
      return true
    },
  }
}
</script>
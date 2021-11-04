<template>
  <!-- contact form start -->
  <section id="contacts-wrap" class="contact-form-wrap light-gray p-t-b-80">
    <b-container>
      <b-row>
        <b-col md="6" class="form-wrap mb-4">
          <div class="section-header mb-5 ">
            <h2 class="font-weight-bold ">Get in Touch With Us !</h2>
            <p>
              Send us an email if you have any feature request, bug reporting, billing or anything.
            </p>
          </div>
          <div v-if="showMessage"
            class="alert alert-success alert-dismissible fade show"
            role="alert"
          >
            <button type="button" class="close" @click="showMessage=false">
              <span aria-hidden="true">&times;</span>
            </button>
            <strong>{{displayedMessage}}</strong>
          </div>

          <form class="row form" role="form">
            <div class="form-group mb-4 col-md-6">
              <label for="">
                <i class="eva eva-person-outline"></i>
              </label>
              <input
                type="text"
                id=""
                v-model="contactForm.name"
                class="form-control"
                placeholder="Your Name"
                required
              />
            </div>
            <div class="form-group mb-4 col-md-6">
              <label for="">
                <i class="eva eva-email-outline"></i>
              </label>
              <input
                type="email"
                required
                v-model="contactForm.email"
                class="form-control"
                placeholder="Your email"
                aria-describedby="helpId"
              />
            </div>
            <div class="form-group mb-4 col-md-12">
              <label for="">
                <i class="eva eva-edit-2-outline"></i>
              </label>
              <input
                type="text"
                required
                v-model="contactForm.subject"
                class="form-control"
                placeholder="subject"
                aria-describedby="helpId"
              />
            </div>

            <div class="form-group mb-4 col-md-12">
              <label for="">
                <i class="eva eva-edit-outline"></i>
              </label>
              <textarea
                class="form-control"
                v-model="contactForm.message"
                placeholder="Please enter message here"
                required
                rows="5"
              ></textarea>
            </div>

            <div class="form-group  col-md-12">
              <button
                type="button"
                class="btn btn-block half-button form-submit-button btn-large btn-gradient"
                @click="sendContactUsEmail"
                :disabled="showMessage"
              >
                Send message
              </button>
            </div>
          </form>
        </b-col>
        <div class="col-md-6 align-items-center">
          <div class="right-contact-wrap ml-5">
            <img
              src="@/assets/images/landing/svg/contact2.svg"
              class="img-responsive zoom-fade"
              alt="Image"
            />

            <!-- <div class="address-wrap mb-4">-->
          </div>
        </div>
      </b-row>
    </b-container>
  </section>

  <!-- contact form end -->
</template>

<script>
  import axios from "axios";
  import {default as API_ENDPOINTS} from "../../api";
  export default {
    data() {
      return {
        contactForm: {
          name: "",
          email: "",
          subject: "",
          message: ""
        },
        displayedMessage: "",
        showMessage: false
      }
    },
    methods: {
      isFormValidate() {
        return this.contactForm.name != "" && this.contactForm.email != "" &&
            this.contactForm.subject != "" && this.contactForm.message != ""
      },
      closeShowMessage() {
          setTimeout(() => {
            this.showMessage = false
          }, 6000)
      },
      sendContactUsEmail() {
        if (!this.isFormValidate()) {
          return
        }
        let data = this.contactForm
        return new Promise((resolve) => {
          axios({url: API_ENDPOINTS.CONTACT_US, data: data, method: 'POST' })
              .then(resp => {
                resolve(resp)
                this.displayedMessage = "Message sent, and we will get back to you shortly!"
                this.showMessage = true
                this.closeShowMessage()
              })
              .catch(() => {
                this.displayedMessage = "Message failed to send, please try again"
                this.showMessage = true
                this.closeShowMessage()
              })
        })
      }
    }
  }
</script>
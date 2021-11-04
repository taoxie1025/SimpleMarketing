<template>
  <section class="faqs-wrap p-t-b-80" id="faqs-wrap">
    <b-container>
      <div class="row"></div>

      <b-row>
        <b-col md="6" sm="12" class="mb-4 mb-lg-0 mb-md-0 mb-4">
          <div class="section-header mb-5">
            <h2 class="font-weight-bold">FAQs</h2>
            <p>

            </p>
          </div>
          <div v-for="(faq, i) in faqs" role="tablist" v-bind:key="i">
            <b-card no-body class="mb-4 o-hidden">
              <b-card-header header-tag="header" class="p-3" role="tab">
                <b-button
                  block
                  href="#"
                  v-b-toggle="`accordion-` + i"
                  variant="bg-transparent"
                  class="font-weight-700 text-14"
                  >{{ faq.question }}</b-button
                >
              </b-card-header>
              <b-collapse
                :id="`accordion-` + i"
                accordion="my-accordion"
                role="tabpanel"
                :visible="faq.visible"
              >
                <b-card-body>
                  <b-card-text>{{ faq.answer }}</b-card-text>
                </b-card-body>
              </b-collapse>
            </b-card>
          </div>
        </b-col>
        <!-- end of col -->
        <b-col md="5" class="offset-md-1 col-sm-12 d-flex align-items-center">
          <img
            src="@/assets/images/landing/svg/faqs2.svg"
            class="img-responsive w-100"
            alt
          />
        </b-col>
      </b-row>
    </b-container>
  </section>
</template>
<script>
export default {
  data() {
    return {
      faqs: [
        {
          "question": "Do I need a credit card to try it out and see if it suits my business?",
          "answer": "Absolutely no. You can upgrade or downgrade your subscription plan at any time and you will be charged on a prorated basis. Our free plan will be free forever and there is no trail period, and your monthly quota will be renewed every 30 days.",
          "visible": false
        },
        {
          "question": "Will I be able to charge my subscribers for premium contents?",
          "answer": "We are working on the feature currently, and it will be available soon."
        },
        {
          "question": "If I have my content already, can I bring them over?",
          "answer": "Yes, we provided a html editor where you can manage your content."
        },
        {
          "question": "Is there a landing page provided to allow my clients to subscribe to my business?",
          "answer": "Yes."
        },
        {
          "question": "Can I implement my own landing page?",
          "answer": "Yes, you can integrate your your own landing page with our APIs. Api doc is available after you sign in, and we will be providing developer support to Pro and Ultra users."
        }
      ]
    };
  }
};
</script>
<style scoped>
</style>

<template>
  <section id="services-wrap" class="services-wrap light-gray">
    <b-container>
      <b-row>
        <b-col md="12" class="section-header mb-5">
          <h2 class="font-weight-bold">Our Services</h2>
          <p>
            Develop your business strategy by connecting to your audience closely and constantly push your contents.
          </p>
        </b-col>

        <div
          class="service-wrap col-md-6 mb-4 col-lg-3 col-sm-12"
          v-for="(service, key) in services"
          :key="key"
          :class="[
            key == 4 || key == 5 || key == 6 || key == 7 ? 'mb-lg-0' : ''
          ]"
        >
          <b-card no-body class="o-hidden ">
            <div class="card-header p-0 text-center ">
              <img class="img-responsive" :src="service.serviceImage" />
            </div>
            <b-card-body>
              <h3 class="card-title text-capitalize font-weight-bold  ">
                {{ service.serviceName }}
              </h3>
              <p class="card-text">{{ service.serviceText }}</p>
            </b-card-body>
          </b-card>
        </div>
      </b-row>
    </b-container>
  </section>
</template>
<script>
export default {
  data() {
    return {
      services: [
        {
          serviceImage: require("@/assets/images/landing/services/services16.svg"),
          serviceName: "Email Campaign",
          serviceText:
              "Automate the delivery of your contents to your subscribers on a schedule of your choice"
        },
        {
          serviceImage: require("@/assets/images/landing/services/services3.svg"),
          serviceName: "Template Editing",
          serviceText:
              "Customize your contents with the simple and intuitive html editor"
        },
        {
          serviceImage: require("@/assets/images/landing/services/services17.svg"),
          serviceName: "Landing Page",
          serviceText:
            "A landing page is provided for you to convert your audience to your subscriber"
        },
        {
          serviceImage: require("@/assets/images/landing/services/services15.svg"),
          serviceName: "API Ready",
          serviceText:
              "Well documented API is ready in case if you want to implement your own landing page"
        }
      ]
    };
  }
};
</script>


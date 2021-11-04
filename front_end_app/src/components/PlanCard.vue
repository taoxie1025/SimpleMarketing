<template>
  <div>
    <div class="pricing-four__card">
      <b-card :class="selectedPlan==plan.value ? 'primary' : ''">
        <div class="">
          <div class="pricing-four__contents">
            <div class="pricing-four__titles text-center">
              <h2 class="font-weight-bold mt-4 mb-4">{{ plan.name }}</h2>
            </div>
            <div class="pricing-four__values  text-center mt-4 mb-4">
              <p class="text-25">
                $ <span class="text-50 font-weight-bold">{{ plan.price }}</span> / month
              </p>
            </div>

            <ul class="list-group list-group-flush mt-5">
              <li class="list-group-item bg-transparent border-0 pb-0">
                <i class="eva eva-checkmark-outline align-middle mr-2"></i>
                <span class="font-weight-800">{{ plan.quota }}</span> emails per month
              </li>
              <li v-for="(feature) in plan.feature" v-bind:key="feature.name" class="list-group-item bg-transparent border-0 pb-0">
                <div v-if="feature.enabled">
                  <i class="eva eva-checkmark-outline align-middle mr-2"></i>
                  <span>{{ feature.name }}</span>
                </div>
                <div v-if="!feature.enabled">
                  <i class=" mr-4"></i>
                  <span class="text-mute"><strike>{{ feature.name }}</strike></span>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </b-card>
    </div>
  </div>
</template>

<script>

export default {
  components: {},
  data() {
    return {
    }
  },
  props:[
    'plan',
    'selectedPlan'
  ],
  methods: {
  },
  computed: {
  },
  mounted() {

  },
  created() {
  }
}
</script>
<style>
</style>
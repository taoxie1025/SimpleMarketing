<template>
  <div class="">
    <dxIntro></dxIntro>
    <dx-feature></dx-feature>
    <dx-service></dx-service>
    <dx-workCarousel></dx-workCarousel>
    <dx-productDetails></dx-productDetails>
    <dx-priceTable></dx-priceTable>
    <dx-team></dx-team>
    <dx-faq></dx-faq>
  </div>
</template>
<script>
import intro from "@/components/intros/introEight";
import feature from "@/components/features/feature";
import service from "@/components/service/service";
import workCarousel from "@/components/works/workCarousel";
import productDetails from "@/components/product-details/best";
import priceTable from "@/components/price/priceThree";
import team from "@/components/team/team";
import faq from "@/components/faq/faq";

export default {
  metaInfo: {
    bodyAttrs: {
      class: ["landing-indigo-deep"]
    }
  },
  components: {
    dxIntro: intro,
    dxFeature: feature,
    dxService: service,
    dxWorkCarousel: workCarousel,
    dxProductDetails: productDetails,
    dxPriceTable: priceTable,
    dxTeam: team,
    dxFaq: faq
  },
  methods: {
    handleStorageChangeEvent (event) {
      if(event.key === 'logged_in') {
        location.reload()
      }
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
  mounted() {
    window.addEventListener('storage', this.handleStorageChangeEvent, false)
  }
};
</script>

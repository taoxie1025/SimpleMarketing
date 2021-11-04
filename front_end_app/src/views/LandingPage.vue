<template>
  <div class="">
    <dxIntro></dxIntro>
  </div>
</template>
<script>
import intro from "@/components/intros/introEleven";

export default {
  metaInfo: {
    bodyAttrs: {
      class: ["landing-gradient-steel-gray"]
    }
  },
  components: {
    dxIntro: intro,
  }
};
</script>

<template>
  <!-- extra features start -->
  <section class="best-wrap p-t-b-80 light-gray" id="extra-feature-wrap">
    <b-container>
      <b-row>
        <b-col md="12" class="section-header mb-5">
          <h2 class="font-weight-bold">{{ bestHeading }}</h2>
          <p>{{ bestTitle }}</p>
        </b-col>

        <b-col md="6"
          class="mb-4"
          v-for="(best, key) in bestContent"
          :key="key"
          :class="[key == 2 || key == 3 ? 'mb-lg-0' : '']"
        >
          <b-card class="o-hidden text-left">
            <div class="">
              <div class="left_icon float-left w-20 p-3 text-center">
                <i class="eva" :class="best.bestIcon"></i>
              </div>
              <div class="right_content w-70 float-left p-3">
                <h4 class="card-title font-weight-bold">
                  {{ best.bestTitle }}
                </h4>
                <p class="card-text">{{ best.bestText }}</p>
              </div>
            </div>
          </b-card>
        </b-col>
        <!-- end of col -->
      </b-row>
    </b-container>
  </section>

  <!-- extra features end -->
</template>
<script>
export default {
  data() {
    return {
      bestHeading: "Applications",
      bestTitle:
        "No matter you are a content creator or entrepreneur, we've got your back.",
      bestContent: [
        {
          bestTitle: "Youtuber",
          bestText:
            "Keep your premium subscribers informed with your paid contents as an episode.",
          bestIcon: ""
        },
        {
          bestTitle: "Writer",
          bestText:
            "Push out your new book to your subscribers section by section on a schedule of your choice.",
          bestIcon: ""
        },
        {
          bestTitle: "Software Developer",
          bestText:
            "Share your learning and experience with your subscribers and build up your influence in the field.",
          bestIcon: ""
        },
        {
          bestTitle: "Business Owner",
          bestText:
            "Advertise your products or services to your subscribers on a timely manner.",
          bestIcon: ""
        }
      ]
    };
  }
};
</script>

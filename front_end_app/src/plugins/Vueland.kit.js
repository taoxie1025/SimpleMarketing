import BootstrapVue from "bootstrap-vue";
import VueLazyload from "vue-lazyload";
import VueMeta from "vue-meta";

import "@/assets/styles/sass/themes/lite-purple.scss";

// locale.use(lang);

export default {
    install(Vue) {
        Vue.use(BootstrapVue);
        Vue.use(VueMeta, {
            keyName: "metaInfo",
            attribute: "data-vue-meta",
            ssrAttribute: "data-vue-meta-server-rendered",
            tagIDKeyName: "vmid",
            refreshOnceOnNavigation: true
        });
        Vue.use(VueLazyload, {
            observer: true,
            // optional
            observerOptions: {
                rootMargin: "0px",
                threshold: 0.1
            }
        });
    }
};

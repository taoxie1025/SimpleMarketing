import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import store from './store'
import vuetify from './plugins/vuetify';
import 'material-design-icons-iconfont/dist/material-design-icons.css'
import Axios from 'axios'

// tiptapVuetify plugin
import './plugins/tiptap-vuetify'
import './plugins/gtag'
import '@babel/polyfill'

// dexam template plugins
import VueLand from "./plugins/Vueland.kit";
// import on your project (less then 1KB gziped)
import vueSmoothScroll from "vue2-smooth-scroll";
import infiniteScroll from 'vue-infinite-scroll'

Vue.use(vueSmoothScroll);
Vue.use(VueLand);
Vue.use(infiniteScroll)

Vue.config.productionTip = false
Vue.prototype.$http = Axios;

const token = localStorage.getItem('token')
if (token) {
  Vue.prototype.$http.defaults.headers.common['Authorization'] = token
}

new Vue({
  router,
  store,
  vuetify,
  VueLand,
  render: h => h(App)
}).$mount('#app')

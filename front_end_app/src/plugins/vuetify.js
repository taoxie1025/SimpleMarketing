import '@mdi/font/css/materialdesignicons.css'
import Vue from 'vue'
import Vuetify from 'vuetify'


Vue.use(Vuetify)
export default new Vuetify({
    iconfont: 'md',
    theme: {
        primary: '#3649d4',
        success: '#3cd1c2',
        info: '#ffaa2c',
        error: '#f83e70',
        completes: '#3cd1c2',
        ongoings: '#ffaa2c',
        overdues: '#f83e70'
    }
})
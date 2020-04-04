import Vue from 'vue'
import App from './App.vue'
Vue.config.productionTip = false

import 'vue-material/dist/vue-material.min.css'
import 'vue-material/dist/theme/default.css'
import { Datetime } from 'vue-datetime'
import 'vue-datetime/dist/vue-datetime.css'
import { ContainerMixin, ElementMixin } from 'vue-slicksort';
import VueMaterial from 'vue-material'

Vue.use(VueMaterial)
Vue.component('datetime', Datetime)

new Vue({
        render: h => h(App),
}).$mount('#app')

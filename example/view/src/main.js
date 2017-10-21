// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'

import Mint from 'mint-ui'
import 'mint-ui/lib/style.css'
import VueAMap from 'vue-amap'

Vue.use(Mint)
Vue.use(VueAMap)
VueAMap.initAMapApiLoader({
  key: '8cc3327f86ecdbc7840b2bee9e0997da',
  plugin: ['Geolocation','Autocomplete','Scale','OverView','ToolBar']
});

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})

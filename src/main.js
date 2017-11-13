import Vue from 'vue'
import VueRouter from 'vue-router'

import App from './App.vue'
import routes from './utils/routes'
import store from './utils/store'

import ElementUI from 'element-ui'
import locale from 'element-ui/lib/locale/lang/en'
import 'element-ui/lib/theme-default/index.css'

Vue.use(VueRouter)
Vue.use(ElementUI, {locale})

const router = new VueRouter({
    mode: 'hash', // hash/history
    routes: routes,
    linkActiveClass: '',
    linkExactActiveClass: 'current',
    scrollBehavior (to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { x: 0, y: 0 }
        }
    },
})

Vue.filter('fltUpperCase', function (value) {
    if (typeof(value) === 'string') {
        return value.toLocaleUpperCase()
    }
    return value.toString().toUpperCase()
})

new Vue({
    el: '#app',
    store,
    router,
    render: h => h(App),
})
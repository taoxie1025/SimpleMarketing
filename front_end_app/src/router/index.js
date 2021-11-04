import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store/index.js'
import Dashboard from '../views/Dashboard.vue'
import Articles from '../views/Articles.vue'
import Home from '../views/Home.vue'
import LandingPage from '../views/LandingPage.vue'
import ApiDoc from '../views/ApiDoc'
import Subscribers from '../views/Subscribers.vue'
import Plan from '../views/Plan.vue'
import Account from '../views/Account.vue'
import Recovery from '../views/Recovery.vue'
import Unsubscribe from '../views/Unsubscribe.vue'
import Support from '../views/Support.vue'
import Admin from '../views/Admin.vue'

Vue.use(VueRouter)

  const routes = [
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: '/project/:projectId',
      name: 'Articles',
      component: Articles,
      meta: {
        requiresAuth: true
      },
      props: true
    },
    {
      path: "/",
      name: 'Home',
      component: Home
    },
    {
      path: "/recovery",
      name: 'Recovery',
      component: Recovery
    },
    {
      path: "/unsubscribe",
      name: 'Unsubscribe',
      component: Unsubscribe
    },
    {
      path: "/p/:projectIdShort",
      name: 'LandingPage',
      component: LandingPage
    },
    {
      path: "/documents/:projectId/api",
      name: 'ApiDoc',
      component: ApiDoc
    },
    {
      path: '/subscribers/:projectId',
      name: 'Subscribers',
      component: Subscribers,
      meta: {
        requiresAuth: true
      },
      props: true
    },
    {
      path: '/plan',
      name: 'Plan',
      component: Plan,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: '/account',
      name: 'Account',
      component: Account,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: '/support',
      name: 'Support',
      component: Support,
      meta: {
        requiresAuth: true
      }
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin,
      meta: {
        requiresAuth: true
      }
    },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  if(to.matched.some(record => record.meta.requiresAuth)) {
    if (store.getters.isLoggedIn) {
      next()
      return
    }
    next('/')
  } else if (to.name == 'Auth' && store.getters.isLoggedIn && store.getters.email){
      next('/')
      return
  } else {
    next()
  }
})
export default router

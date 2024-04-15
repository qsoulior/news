import { createRouter, createWebHistory } from "vue-router"

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/:page?",
      name: "list",
      props: (route) => {
        const queries = route.query.page
        const queryRaw = Array.isArray(queries) ? queries[0] : queries
        if (queryRaw == null) return { page: undefined }
        const query = parseInt(queryRaw)
        return { page: isNaN(query) ? undefined : query }
      },
      component: () => import("@/views/ListView.vue")
    },
    {
      path: "/news/:id",
      name: "item",
      props: true,
      component: () => import("@/views/ItemView.vue")
    }
  ]
})

export default router

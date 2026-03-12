import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      children: [
        { path: '', name: 'home', component: () => import('@/views/Home.vue') },
        { path: 'explore', name: 'explore', component: () => import('@/views/Explore.vue') },
        { path: 'login', name: 'login', component: () => import('@/views/Login.vue') },
        { path: 'register', name: 'register', component: () => import('@/views/Register.vue') },
        { path: 'settings/profile', name: 'settings-profile', component: () => import('@/views/settings/ProfileSettings.vue') },
        { path: 'settings/tokens', name: 'settings-tokens', component: () => import('@/views/settings/TokenSettings.vue') },
        {
          path: 'admin',
          component: () => import('@/views/admin/AdminLayout.vue'),
          children: [
            { path: '', name: 'admin', component: () => import('@/views/admin/AdminDashboard.vue') },
            { path: 'users', name: 'admin-users', component: () => import('@/views/admin/AdminUsers.vue') },
            { path: 'looms', name: 'admin-looms', component: () => import('@/views/admin/AdminLooms.vue') },
          ],
        },
        {
          path: ':owner/:loom',
          component: () => import('@/views/loom/LoomLayout.vue'),
          children: [
            { path: '', name: 'loom', component: () => import('@/views/loom/LoomOverview.vue') },
            { path: 'tree', name: 'loom-tree', component: () => import('@/views/loom/EntityTree.vue') },
            { path: 'tree/:treePath(.*)', name: 'loom-tree-path', component: () => import('@/views/loom/EntityTree.vue') },
            { path: 'blob/:path(.*)', name: 'loom-blob', component: () => import('@/views/loom/EntityViewer.vue') },
            { path: 'log', name: 'loom-log', component: () => import('@/views/loom/CheckpointLog.vue') },
            { path: 'streams', name: 'loom-streams', component: () => import('@/views/loom/Streams.vue') },
            { path: 'weave-requests', name: 'loom-wrs', component: () => import('@/views/loom/WRList.vue') },
            { path: 'weave-requests/new', name: 'loom-wr-new', component: () => import('@/views/loom/WRNew.vue') },
            { path: 'weave-requests/:number', name: 'loom-wr', component: () => import('@/views/loom/WRDetail.vue') },
            { path: 'settings', name: 'loom-settings', component: () => import('@/views/loom/Settings.vue') },
            { path: 'settings/collaborators', name: 'loom-settings-collaborators', component: () => import('@/views/loom/SettingsCollaborators.vue') },
            { path: 'settings/webhooks', name: 'loom-settings-webhooks', component: () => import('@/views/loom/SettingsWebhooks.vue') },
            { path: 'settings/labels', name: 'loom-settings-labels', component: () => import('@/views/loom/SettingsLabels.vue') },
          ],
        },
        { path: ':owner', name: 'profile', component: () => import('@/views/Profile.vue') },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      component: () => import('@/layouts/DefaultLayout.vue'),
      children: [
        { path: '', name: 'not-found', component: () => import('@/views/NotFound.vue') },
      ],
    },
  ],
})

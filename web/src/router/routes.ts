import { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ path: '', component: () => import('pages/IndexPage.vue') }],
  },
  {
    path: '/signup',
    component: () => import('layouts/AuthLayout.vue'),
    children: [{ path: '', component: () => import('pages/SignupPage.vue') }],
  },
  {
    path: '/sandbox',
    component: () => import('layouts/SandboxLayout.vue'),
    children: [{ path: '', component: () => import('pages/SignupPage.vue') }],
  },
  {
    path: '/privacy',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ path: '', component: () => import('pages/PrivacyPage.vue') }],
  },
  {
    path: '/auth',
    component: () => import('layouts/AuthLayout.vue'),
    children: [
      {
        path: 'confirm/:token',
        component: () => import('pages/auth/ConfirmPage.vue'),
      },
      {
        path: 'login',
        component: () => import('pages/auth/LoginPage.vue'),
      },
      {
        path: 'reset-password/:token',
        component: () => import('pages/auth/ResetPasswordPage.vue'),
      },
      {
        path: 'forgotpassword',
        component: () => import('pages/auth/LostPasswordPage.vue'),
      },
      {
        path: 'mfa/:token',
        component: () => import('pages/auth/MFAPage.vue'),
      },
      {
        path: 'otp',
        component: () => import('src/pages/auth/OTPLoginPage.vue'),
      },
      {
        path: 'otp/:token',
        component: () => import('src/pages/auth/OTPVerifyPage.vue'),
      },
    ],
  },
  {
    path: '/oauth/:type/:provider',
    component: () => import('pages/auth/OAuthPage.vue'),
  },
  {
    path: '/user',
    component: () => import('layouts/UserLayout.vue'),
    children: [
      {
        path: 'dashboard',
        component: () => import('pages/user/DashboardPage.vue'),
      },
      {
        path: 'account',
        component: () => import('pages/user/AccountPage.vue'),
      },
      {
        path: 'password',
        component: () => import('pages/user/PasswordPage.vue'),
      },
      {
        path: 'settings',
        component: () => import('pages/user/SettingsPage.vue'),
      },
      {
        path: 'profile-image',
        component: () => import('pages/user/ProfileImagePage.vue'),
      },
    ],
  },
  {
    path: '/examples',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ path: '', component: () => import('pages/ExamplePage.vue') }],
  },
  {
    path: '/admin',
    component: () => import('layouts/UserLayout.vue'),
    children: [
      {
        path: 'organization',
        children: [
          {
            path: 'dashboard',
            component: () =>
              import('pages/admin/organization/DashboardPage.vue'),
          },
          {
            path: 'information',
            component: () =>
              import('pages/admin/organization/InformationPage.vue'),
          },
          {
            path: 'users',
            component: () => import('pages/admin/organization/UsersPage.vue'),
          },
          {
            path: 'forms',
            component: () => import('pages/admin/organization/FormsPage.vue'),
          },
          {
            path: 'tree/:slug',
            component: () => import('pages/admin/organization/TreePage.vue'),
          },
          {
            path: 'trees',
            component: () => import('pages/admin/organization/TreesPage.vue'),
          },
        ],
      },
      {
        path: 'content',
        children: [
          {
            path: 'setup',
            component: () =>
              import('pages/admin/organization/ContentSetupPage.vue'),
          },
          {
            path: 'manager',
            component: () =>
              import('pages/admin/organization/ContentManagerPage.vue'),
          },
        ],
      },
      {
        path: 'cms',
        component: () => import('pages/admin/organization/CMSPage.vue'),
      },
      {
        path: 'board/:slug',
        component: () => import('pages/admin/organization/BoardPage.vue'),
      },
    ],
  },
  {
    path: '/c',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: 'b/:slug', component: () => import('pages/cms/BoardPage.vue') },
      { path: 'p/:slug', component: () => import('pages/cms/PostPage.vue') },
    ],
  },
  {
    path: '/sandbox',
    component: () => import('layouts/SandboxLayout.vue'),
    children: [
      {
        path: 'tiptap',
        component: () => import('pages/sandbox/TiptapPage.vue'),
      },
      {
        path: 'advanced-tiptap',
        component: () => import('pages/sandbox/AdvancedTiptapPage.vue'),
      },
      {
        path: 'api',
        component: () => import('pages/sandbox/ApiPage.vue'),
      },
    ],
  },
  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
];

export default routes;

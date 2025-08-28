import { boot } from 'quasar/wrappers';
import axios, { AxiosInstance } from 'axios';
import { useUserStore } from 'src/stores/user-store';
import { v4 as uuidv4 } from 'uuid';

const accessToken = useUserStore();

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $axios: AxiosInstance;
  }
}

// Be careful when using SSR for cross-request state pollution
// due to creating a Singleton instance here;
// If any client changes this (global) instance, it might be a
// good idea to move this instance creation inside of the
// "export default () => {}" function below (which runs individually
// for each client)

const apiPort = process.env.API_PORT ? `:${process.env.API_PORT}` : '';
const apiDomain =
  window.location.protocol + '//' + window.location.hostname + apiPort;

const instance = axios.create({
  baseURL: apiDomain,
  withCredentials: true,
});

const api = () => {
  const tokenHeader = accessToken.getHeader;
  if (tokenHeader) {
    instance.defaults.headers.common['Authorization'] = tokenHeader;
  } else {
    instance.defaults.headers.common['Authorization'] = '';
  }
  //for istio
  instance.defaults.headers.common['x-request-id'] = uuidv4();
  instance.defaults.headers.common['x-host'] =
    window.location.host.split(':')[0];
  return instance;
};

export default boot(({ app }) => {
  // for use inside Vue files (Options API) through this.$axios and this.$api

  app.config.globalProperties.$axios = axios;
  // ^ ^ ^ this will allow you to use this.$axios (for Vue Options API form)
  //       so you won't necessarily have to import axios in each vue file

  app.config.globalProperties.$api = api;
  // ^ ^ ^ this will allow you to use this.$api (for Vue Options API form)
  //       so you can easily perform requests against your app's API
});

export { api };

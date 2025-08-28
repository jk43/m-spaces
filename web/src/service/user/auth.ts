import { api } from 'src/boot/axios';
import { useUserStore } from 'src/stores/user-store';
import { LoginSuccessResponse } from 'src/types';
import { useOrgStore } from 'src/stores/org-store';

const storeOrg = useOrgStore();
const storeUser = useUserStore();
let timer: ReturnType<typeof setInterval> | null = null;

export async function refreshToken(runLoop: boolean, router: any) {
  try {
    const result = await api().get('/user/refreshtoken', {});
    storeUser.setUserInfo(result.data.data as LoginSuccessResponse);
    if (runLoop) {
      refreshTokenInterval(8, router);
    }
  } catch (error) {
    if (timer !== null) {
      clearInterval(timer);
    }
    storeUser.setUserInfo(null);
    if (router !== null) {
      console.log('storeOrg.active: ', storeOrg.active);
      if (!storeOrg.active) {
        return;
      }
      if (
        !window.location.pathname.startsWith('/auth/') &&
        !window.location.pathname.startsWith('/oauth/') &&
        !window.location.pathname.startsWith('/c/') &&
        window.location.pathname !== '/signup' &&
        window.location.pathname !== '/examples/'
      ) {
        alert('Session expired, please login again');
        window.location.href = '/auth/login';
      }
    }
  }
}

export function refreshTokenInterval(minute: number, router: any) {
  timer = setInterval(() => {
    refreshToken(false, router);
    console.log('refreshTokenInterval: ');
  }, 1000 * 60 * minute);
}

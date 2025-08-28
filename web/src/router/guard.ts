import { useOrgStore } from 'src/stores/org-store';
import { useUserStore } from 'src/stores/user-store';

export const beforeEach = async (to, from, next) => {
  let res = checkOrganizationStatus(to, from);
  if (res) {
    next(res);
    return;
  }
  res = checkLoginStatus(to, from);
  if (res) {
    next(res);
    return;
  }
  next();
};

const checkOrganizationStatus = (to, from): string | null => {
  const storeOrg = useOrgStore();
  if (storeOrg.active == false) {
    if (to.path == '/pagenotfound') {
      return null;
    } else {
      return '/pagenotfound';
    }
  } else {
    return null;
  }
};

const checkLoginStatus = (to, from): string | null => {
  const storeUser = useUserStore();
  if (to.path == '/auth/login' && storeUser.isLoggedIn) {
    return '/user/dashboard';
  }
  if (storeUser.isLoggedIn) {
    return null;
  } else {
    if (to.path == '/auth/login') {
      return null;
    }
    const firstPath = to.path.split('/')[1];
    if (firstPath !== 'user') {
      return null;
    }
    return '/auth/login';
  }
};

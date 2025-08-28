import { defineStore } from 'pinia';
import { LoginSuccessResponse } from 'src/types';
import { useOrgStore } from 'src/stores/org-store';

export const useUserStore = defineStore('user', {
  state: () => ({
    accessToken: '',
    isLoggedIn: false,
    firstName: '',
    lastName: '',
    profileImage: '',
    email: '',
    metadata: {},
    oauth: false,
    id: '',
  }),
  actions: {
    setUserInfo(data: LoginSuccessResponse | null) {
      if (data === null) {
        this.accessToken = '';
        this.firstName = '';
        this.lastName = '';
        this.profileImage = '';
        this.email = '';
        this.isLoggedIn = false;
        this.metadata = {};
        this.id = '';
        return;
      }
      if (data.accessToken) {
        this.isLoggedIn = true;
      } else {
        this.isLoggedIn = false;
      }
      this.accessToken = data.accessToken;
      this.firstName = data.info.firstName;
      this.lastName = data.info.lastName;
      this.profileImage = data.info.profileImage;
      this.email = data.info.email;
      this.metadata = data.info.metadata;
      this.id = data.info.id;
      this.oauth = !['admin', 'self'].includes(data.info.registerMethod);
    },
  },
  getters: {
    getHeader: (state) => {
      if (!state.isLoggedIn) {
        return null;
      }
      return 'Bearer ' + state.accessToken;
    },
    getProfileImage: (state) => {
      if (state.profileImage === '') {
        return '';
      }
      if (state.profileImage.startsWith('https://')) {
        return state.profileImage;
      }
      return useOrgStore().cdnAddr + '/' + state.profileImage;
    },
    getId: (state) => state.id,
  },
});

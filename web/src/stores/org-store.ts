import { defineStore } from 'pinia';
import {
  OrganizationSettingResponse,
  OrganizationItems,
  FormInput,
  LabelValuePair,
} from 'src/types';
import { capitalize } from 'vue';
interface Org {
  active: boolean;
  host: string;
  cdnAddr: string;
  name: string;
  allowSelfMemberRegistration: boolean;
  settings: {
    forms: { [key: string]: FormInput };
    roles: string[];
  };
  items: OrganizationItems;
  itemReady: true;
  otp: boolean;
}

export const useOrgStore = defineStore('org', {
  state: (): Org => ({
    active: true,
    host: '',
    cdnAddr: '',
    name: '',
    allowSelfMemberRegistration: false,
    settings: {
      forms: {},
      roles: [] as LabelValuePair[],
    },
    items: {},
    itemReady: false,
    otp: false,
  }),
  getters: {
    getForm: (state) => (key: string) => {
      return state.settings.forms[key];
    },
    getRoles: (state) => () => {
      return state.settings.roles;
    },
  },
  actions: {
    async setSettings(data: OrganizationSettingResponse) {
      this.active = data.active;
      this.host = data.host;
      this.cdnAddr = data.cdnAddr;
      this.name = data.name;
      this.allowSelfMemberRegistration = data.allowSelfMemberRegistration;
      this.otp = data.otp;
      //this.settings = data.settings;
      for (const [key, val] of Object.entries(data.settings.forms)) {
        this.settings.forms[key] = val.sort((a, b) => a.order - b.order);
        // for (const fi of val) {
        //   //val.sort((a, b) => a.order - b.order);
        // }
      }
      data.settings.roles
        .sort((a, b) => a.order - b.order)
        .forEach((r) => {
          this.settings.roles.push({
            label: capitalize(r.role),
            value: r.role,
          });
        });
    },
    async setItems(items: OrganizationItems) {
      for (const key in items) {
        this.items[key] = items[key].sort((a, b) => a.order - b.order);
      }
      this.itemReady = true;
    },
  },
});

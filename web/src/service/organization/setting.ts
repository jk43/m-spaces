import { api } from 'src/boot/axios';
import { useOrgStore } from 'src/stores/org-store';
import { OrganizationSettingResponse, Response } from 'src/types';

const storeOrg = useOrgStore();
const timer: ReturnType<typeof setInterval> | null = null;

export async function getOrganizationSettings() {
  try {
    const result = await api().get<Response<OrganizationSettingResponse>>(
      '/organization/settings',
      {}
    );
    storeOrg.setSettings(result.data.data);
  } catch (error) {
    storeOrg.setSettings({ active: false } as OrganizationSettingResponse);
  }
}

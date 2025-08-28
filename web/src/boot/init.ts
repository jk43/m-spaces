import { boot } from 'quasar/wrappers';
import { refreshToken } from 'src/service/user/auth';
import { getOrganizationSettings } from 'src/service/organization/setting';

// Tiptap CSS 전역 import
import '../css/tiptap-content.css';

// "async" is optional;
// more info on params: https://v2.quasar.dev/quasar-cli/boot-files
export default boot(async ({ app, router }) => {
  await refreshToken(true, router);
  await getOrganizationSettings();
});

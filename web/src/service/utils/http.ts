import { Ref } from 'vue';
import { Router, LocationQuery } from 'vue-router';
import { useUserStore } from 'src/stores/user-store';

const accessToken = useUserStore().getHeader;

export const queryBuilder = (
  router: Router,
  params: Ref<{ [key: string]: any }>
): string => {
  const query = new URLSearchParams(params.value).toString();
  router.push({ path: router.currentRoute.value.path, query: params.value });
  return '?' + query;
};

interface PaginationAndOtherReturnType {
  pagination: { [key: string]: any };
  others: { [key: string]: any };
}

export const separatePagination = (obj: {
  [key: string]: any;
}): PaginationAndOtherReturnType => {
  const output = { pagination: {}, others: {} };
  for (const [k, v] of Object.entries(obj)) {
    if (quasarPaginationVars.includes(k)) {
      output.pagination[k] = v;
    } else {
      output.others[k] = v;
    }
  }
  return output;
};

export const quasarPaginationVars: Array<string> = [
  'sortBy',
  'descending',
  'page',
  'rowsPerPage',
  'rowsNumber',
];

export const GetWebSocket = (channel: string) => {
  let token = '';
  if (accessToken !== '' || accessToken !== null) {
    token = accessToken?.replace('Bearer ', '');
  }
  return new WebSocket('ws://localhost:8081/ws/' + channel + '?t=' + token);
};

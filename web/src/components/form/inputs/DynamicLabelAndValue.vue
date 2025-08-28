<template>
  <transition name="fade">
    <div
      v-show="settings.show.showLabelValue || settings.show.showOption"
      class="q-pa-xs rounded-borders bg-grey-3"
    >
      <div class="flex flex-center text-body1 q-py-sm text-grey-8">
        {{ data.name }} <q-space />
        <q-btn @click="add" flat dense><q-icon name="add" /></q-btn>
      </div>
      <div v-if="settings.show.showOption">
        <div
          class="row q-mb-xs"
          v-for="(item, index) in data.value.options"
          :key="index"
          horizontal
        >
          <div class="col-1 flex flex-center">{{ index + 1 }}</div>
          <div class="col-10">
            <q-input
              lazy-rules
              filled
              :error="error.error"
              :errorMessage="error.message"
              label="Value"
              v-model="data.value.options[index]"
              hide-bottom-space
              width="50%"
              dense
            />
          </div>
          <div class="col-1 flex flex-center">
            <q-btn flat dense
              ><q-icon name="remove" @click="deleteOption(index)"
            /></q-btn>
          </div>
        </div>
      </div>
      <div v-if="settings.show.showLabelValue">
        <div
          v-for="(item, index) in data.value.labelValue"
          :key="index"
          class="row q-gutter-xs q-mb-xs"
        >
          <div class="col-1 flex flex-center">{{ index + 1 }}</div>
          <div class="col">
            <q-input
              lazy-rules
              filled
              :error="error.error"
              :errorMessage="error.message"
              label="Label"
              v-model="item.label"
              hide-bottom-space
              width="50%"
              dense
            />
          </div>
          <div class="col">
            <q-input
              lazy-rules
              filled
              :error="error.error"
              :errorMessage="error.message"
              label="Value"
              v-model="item.value"
              hide-bottom-space
              width="50%"
              dense
            />
          </div>
          <div class="col-1 flex flex-center">
            <q-btn flat dense
              ><q-icon name="remove" @click="deleteValueLabel(index)"
            /></q-btn>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>
<script setup>
import { reactive, defineProps, ref } from 'vue';
const props = defineProps(['data', 'settings', 'error']);
const data = reactive(props.data);
const error = reactive(props.error);
const settings = reactive(props.settings);
const add = () => {
  if (settings.show.showOption) {
    if (!data.value.options) {
      data.value.options = [];
    }
    data.value.options.push('');
  } else if (settings.show.showLabelValue) {
    if (!data.value.options) {
      data.value.labelValue = [{ label: '', value: '' }];
    }
    data.value.labelValue.push({ label: '', value: '' });
  }
};
const deleteOption = (i) => {
  data.value.options.splice(i, 1);
};
const deleteValueLabel = (i) => {
  data.value.labelValue.splice(i, 1);
};

settings.callback(settings, data);
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity s;
}
.fade-enter,
.fade-leave-to {
  opacity: 0;
}
</style>

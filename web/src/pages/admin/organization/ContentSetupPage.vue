<template>
  <div>
    content setup
    <draggable v-model="node" item-key="id" @end="onEnd">
      <template #item="{ element, index }">
        <q-card
          class="my-card shadow-1 q-mb-none q-mt-none"
          style="min-width: 500px"
        >
          <q-card-section class="ext-white">
            <div class="text-h7">
              {{ element.attributes.label }} {{ element }}
              <q-btn
                icon="add"
                @click="showForm(element, index, 'add')"
                size="8px"
                dense
                flat
              />
              <q-btn
                icon="edit"
                @click="showForm(element, index, 'edit')"
                size="8px"
                dense
                flat
              />
              <q-btn
                icon="delete"
                @click="deleteNode(element, index)"
                size="8px"
                dense
                flat
              />
            </div>
            <div v-if="element.attributes.description" class="text-h8">
              {{ element.attributes.description }}
            </div>
            <!-- 재귀적으로 자식 노드 렌더링 -->
            <TreeComponent
              v-if="element.children && element.children.length"
              :node="element.children"
              @showForm="showForm"
              :onReorder="props.onReorder"
              :deleteHandler="props.deleteHandler"
            ></TreeComponent>
          </q-card-section>
        </q-card>
      </template>
    </draggable>
  </div>
</template>
<script setup>
import draggable from 'vuedraggable';

import TreeComponent from 'src/components/TreeComponent.vue';

import { ref, defineProps, defineEmits } from 'vue';
const props = defineProps({
  node: Object,
  onReorder: Function,
  deleteHandler: Function,
});
const emit = defineEmits(['showForm', 'onReorder']);

const node = ref(props.node);

const deleteNode = (n, index) => {
  props.deleteHandler(node.value, index);
  //node.value.splice(index, 1);
};

const showForm = (n, index, mode) => {
  emit('showForm', n, index, mode);
};

const editNode = (n, index) => {
  props.editHandler(n, index);
};

const onEnd = (event) => {
  if (event.newIndex === event.oldIndex) {
    return;
  }
  props.onReorder(node.value);
};
</script>

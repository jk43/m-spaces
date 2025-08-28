<template>
  <ol class="wtree">
    <draggable v-model="node" item-key="id" @end="onEnd">
      <template #item="{ element, index }">
        <li :class="element.parent_id ? 'child' : 'root'">
          <span class="q-pa-xs rounded-borders">
            <div class="bg-white q-pa-xs">
              <div
                :class="[
                  element.parent_id
                    ? 'text-subtitle2'
                    : 'text-h6 text-weight-bold',
                  'flex',
                  'flex-left',
                ]"
              >
                {{ element.attributes.label }}
                <q-space />
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
              <div v-if="element.attributes.description" class="text-caption">
                {{ element.attributes.description }}
              </div>
            </div>
          </span>
          <TreeComponent
            v-if="element.children && element.children.length"
            :node="element.children"
            @showForm="showForm"
            :onReorder="props.onReorder"
            :deleteHandler="props.deleteHandler"
          ></TreeComponent>
        </li>
      </template>
    </draggable>
  </ol>
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

const toggle = (node) => {
  node.isVisible = !node.isVisible;
};

const emit = defineEmits(['showForm', 'onReorder']);

const node = ref(props.node);

const deleteNode = (n, index) => {
  props.deleteHandler(node.value, index);
};

const showForm = (n, index, mode) => {
  emit('showForm', n, index, mode);
};

const onEnd = (event) => {
  if (event.newIndex === event.oldIndex) {
    return;
  }
  props.onReorder(node.value);
};
</script>

<style lang="scss" scoped>
$border: #ddd;
$border-hover: rgb(71, 71, 71);
$bg-hover: rgba(0, 0, 0, 0.1);
$text: rgba(102, 102, 102, 1);
$text-hover: #000;
$ident: 15px;
$left: -($ident);
$first: #ddf3fe;
$second: #ddebc8;
$third: #fefcd5;
$fourth: #fdd2d2;
$fifth: #f9e2f2;
$sixth: #6a0dad;
$seventh: #7b68ee;
$eighth: #8a2be2;
$ninth: #9370db;
$tenth: #9932cc;
$eleventh: #9a32cd;
$twelfth: #a020f0;
$thirteenth: #a52a2a;
$fourteenth: #add8e6;
$fifteenth: #adff2f;
$sixteenth: #afeeee;
$seventeenth: #b0c4de;
$eighteenth: #b0e0e6;
$nineteenth: #b22222;
$twentieth: #b8860b;

ol {
  margin-left: $ident;
  padding-inline-start: 0px;
  margin-block-start: 0em;
}

.wtree {
  li {
    list-style-type: none;
    margin: 10px 0 10px 5px;
    position: relative;

    &.child:before {
      content: '';
      counter-increment: item;
      position: absolute;
      top: -9px;
      left: $left;
      border-left: 1px solid $border;
      border-bottom: 1px solid $border;
      width: $ident;
      height: 15px;
    }

    &:after {
      position: absolute;
      content: '';
      top: 5px;
      left: $left;
      border-left: 1px solid $border;
      border-top: 1px solid $border;
      width: $ident;
      height: 100%;
    }

    &:last-child:after {
      display: none;
    }

    span {
      display: block;
      //border-radius: 5px;
      //border: 1px solid $border;
      //padding: 10px;
      //color: $text;
      //text-decoration: none;
    }
  }
}

.wtree {
  li {
    span {
      &:hover,
      &:focus {
        color: $text-hover;
        border: 1px solid $border-hover;

        & + ol {
          li {
            span {
              color: $text-hover;
              border: 1px solid $border-hover;
            }
          }
        }
      }

      &:hover + ol,
      &:focus + ol {
        li {
          &:after,
          &:before {
            border-color: $border-hover;
          }
        }
      }
    }
  }
}

li span {
  background-color: $first;
}

li li span {
  background-color: $second;
}

li li li span {
  background-color: $third;
}

li li li li span {
  background-color: $fourth;
}

li li li li li span {
  background-color: $fifth;
}

li li li li li li span {
  background-color: $sixth;
}

li li li li li li li span {
  background-color: $seventh;
}

li li li li li li li li span {
  background-color: $eighth;
}

li li li li li li li li li span {
  background-color: $ninth;
}

li li li li li li li li li li span {
  background-color: $tenth;
}
</style>

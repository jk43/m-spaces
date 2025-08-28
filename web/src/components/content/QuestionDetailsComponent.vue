<template>
  <q-page>
    <div class="q-gutter-md">
      <div class="text-h6 text-weight-bold">ID: {{ data._id }}</div>
      <div class="row">
        <div class="column q-gutter-xs">
          <div class="col text-h6">Slugs</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.slugs }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-xs">
          <div class="col text-h6">Model</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.llm }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-xs">
          <div class="col text-h6">Active</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.is_active }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-xs">
          <div class="col text-h6">Question</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.question }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-xs">
          <div class="col text-h6">Question</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            <Math :text="data.question"></Math>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-xs">
          <div class="col text-h6">Equation</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            <Math :text="data.equation"></Math>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Answer Key</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.answer_key }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Answer Value</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            <Math :text="data.answer_val"></Math>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Choices</div>
          <div
            class="col bg-grey-3 q-pa-xs rounded-borders"
            v-for="(v, k) in data.choices"
            :key="k"
          >
            {{ k + 1 }} : <Math :text="v"></Math>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Diagram</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.diagram }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Diagram Image</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            <img
              :src="`https://cdn-dev.hotdev.com/${data.diagram_image}`"
              alt="Diagram"
            />
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Description</div>
          <div
            class="col bg-grey-3 q-pa-xs rounded-borders"
            v-for="(v, k) in data.desc"
            :key="k"
          >
            {{ k + 1 }} : <Math :text="v"></Math>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Created At</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.created_at }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Created by</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.created_by }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Claim Total</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.claim_total }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Claims</div>
          <div
            class="col bg-grey-3 q-pa-xs rounded-borders"
            v-for="(v, k) in data.claims"
            :key="k"
          >
            {{ k + 1 }} : {{ v }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Rate Count</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.rate_count }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Rate Total</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.rate_total }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Error</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.error }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Error Result</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.error_result_from_ai }}
          </div>
        </div>
      </div>

      <div class="row">
        <div class="column q-gutter-md">
          <div class="col text-h6">Error Prompt</div>
          <div class="col bg-grey-3 q-pa-xs rounded-borders">
            {{ data.error_prompt }}
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import Math from 'src/components/content/MathComponent.vue';
import { defineProps } from 'vue';

const props = defineProps(['data']);
//const question = props.data;
</script>

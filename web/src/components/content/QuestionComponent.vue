<template>
  <div class="question q-ma-md" v-if="question">
    <div class="q-mb-xs">
      <q-badge
        class="question-id"
        outline
        color="primary"
        :label="`ID: ` + question.seq_id"
      />
      <q-menu v-model="showing">
        <q-list style="min-width: 100px">
          <q-item clickable v-close-popup>
            <q-item-section avatar>
              <q-icon color="red" name="report_problem" />
            </q-item-section>
            <q-item-section @click="reportProblem = true"
              >Report a problem</q-item-section
            >
          </q-item>
          <q-item clickable v-close-popup v-if="!isExample">
            <q-item-section avatar>
              <q-icon color="primary" name="save_alt" />
            </q-item-section>
            <q-item-section @click="saveQuestion = true"
              >Save this question</q-item-section
            >
          </q-item>
        </q-list>
      </q-menu>
      <q-dialog v-model="reportProblem">
        <q-card style="width: 500px">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">Report a problem.</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>
          <q-card-section>
            <FormBuilder
              :data="reportProblemData"
              :settings="reportProblemInputs"
              :handler="submitProblemHandler"
            />
          </q-card-section>
        </q-card>
      </q-dialog>
      <q-dialog v-model="saveQuestion">
        <q-card style="width: 500px">
          <q-card-section class="row items-center q-pb-none">
            <div class="text-h6">Save to favorite</div>
            <q-space />
            <q-btn icon="close" flat round dense v-close-popup />
          </q-card-section>
          <q-card-section>
            <FormBuilder
              :data="saveQuestionData"
              :settings="saveQuestionInputs"
              :handler="saveQuestionHandler"
            />
          </q-card-section>
        </q-card>
      </q-dialog>
    </div>
    <div class="q-mb-md"><PathComponent :path="question.path" /></div>
    <div class="text-body1 q-mb-md">
      <Math :text="question.question"></Math>
    </div>
    <div class="q-mt-lg q-mb-lg" v-if="question.diagram_image">
      <img :src="`https://cdn-dev.hotdev.com/${question.diagram_image}`" />
    </div>
    <div v-if="incorrect">
      <q-banner rounded inline-actions class="text-white bg-red">
        The selected answer is incorrect.
        <q-btn label="Show Answer" color="primary" @click="showAnswer" />
      </q-banner>
      <div v-if="!isExample">
        Please rate this question:
        <q-rating v-model="rating" size="2em" :max="5" color="primary" />
      </div>
    </div>
    <div v-if="correct">
      <q-banner rounded inline-actions class="text-white bg-blue">
        The selected answer is correct!
      </q-banner>
      <div v-if="!isExample">
        Please rate this question:
        <q-rating v-model="rating" size="2em" :max="5" color="primary" />
      </div>
    </div>
    <div class="text-subtitle2">
      <div v-for="(c, i) in question.choices" :key="i" class="choice-container">
        <q-radio size="sm" v-model="answer" :val="i + 1" />
        <div :id="`choice-` + String(i)"><Math :text="c"></Math></div>
      </div>
    </div>
    <div
      v-if="showDesc"
      class="text-body2 q-mt-md bg-grey-3 q-pa-sm rounded-borders"
    >
      <div v-for="(c, i) in question.desc" :key="i" class="text-body1">
        <q-badge rounded color="red" :label="`Step ` + (i + 1)" />
        <Math :text="c"></Math>
      </div>
    </div>
    <div class="q-mt-lg">
      <q-btn
        v-if="!disableSubmit"
        label="Submit"
        color="primary"
        @click="submit"
        :disable="answer === null"
      />
      <q-btn
        v-if="disableSubmit"
        label="Next"
        color="secondary"
        @click="next"
      />
    </div>
    <div>{{ question.answer_key }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, onBeforeMount, reactive, ref, computed } from 'vue';
import { api } from 'src/boot/axios';

import PathComponent from 'src/components/content/PathComponent.vue';
import Math from 'src/components/content/MathComponent.vue';
import FormBuilder from 'src/components/form/FormBuilderComponent.vue';

import { Response } from 'src/types';

const props = defineProps(['slug', 'id']);
const id = props.id;
const slug = props.slug;

const showing = ref(false);

const questions = ref([]);
// current index
const index = ref(0);
let isExample = false;
const rating = ref(0);

const reportProblem = ref(false);

const reportProblemInputs = reactive([
  {
    key: 'type',
    name: 'Type',
    options: {
      labelValue: [
        { label: 'Answer is not correct', value: 'Answer is not correct' },
        { label: 'Question is ambigurs ', value: 'ambigurs' },
        { label: 'Other', value: 'Other' },
      ],
    },
    rules: ['Required'],
    type: 'Select',
    editable: true,
  },
  {
    key: 'description',
    name: 'Description(optional)',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'TextArea',
    editable: true,
  },
]);

const reportProblemData = ref({
  type: { name: '', value: '' },
  description: { name: '', value: '' },
});

const submitProblemHandler = async () => {
  alert(1);
};

const saveQuestion = ref(false);
const saveQuestionData = ref({
  description: { name: '', value: '' },
});

const saveQuestionInputs = reactive([
  {
    key: 'description',
    name: 'Description(optional)',
    description: '',
    options: [],
    //rules: "[(val) => (val && val.length > 0) || 'First Name cannot be empty']",
    type: 'TextArea',
    editable: true,
  },
]);

const saveQuestionHandler = async () => {
  alert(2);
  // Send the report to the server
  // console.log(selectedData.value);
  // reportProblem.value = false;
};

onBeforeMount(async () => {
  const res: Response<{ [key: string]: any }> = await api().get(
    `/math/question/${slug}/${id}`
  );
  questions.value = res.data.data;
  // Check if the questions are example
  isExample = questions.value[0].fetched_by === 'example';
});

const question = computed(() => {
  return questions.value[index.value];
});

const answer = ref(null);
const incorrect = ref(false);
const correct = ref(false);
const disableSubmit = ref(false);

const submit = () => {
  // Does question has answer key?
  const answerKey = question.value.answer_key;
  if (answerKey) {
    if (Number(answerKey) === answer.value) {
      incorrect.value = false;
      correct.value = true;
    } else {
      incorrect.value = true;
      correct.value = false;
    }
    disableSubmit.value = true;
    return;
  }
  // If not, check if answer is correct over api
};

const showDesc = ref(false);
const prevAnswer = ref<string | null>(null);

const showAnswer = () => {
  const i = 1;
  const elementId = `choice-${i}`;
  const element = document.getElementById(elementId);
  if (element) {
    prevAnswer.value = elementId;
    element.classList.toggle('highlight');
  }
  showDesc.value = true;
};

const next = () => {
  index.value++;
  if (index.value === questions.value.length) {
    if (isExample) {
      exampleDone();
      return;
    }
  }
  rating.value = 0;
  answer.value = null;
  incorrect.value = false;
  correct.value = false;
  disableSubmit.value = false;
  showDesc.value = false;
  if (prevAnswer.value !== null) {
    document.getElementById(prevAnswer.value)?.classList.remove('highlight');
  }
};

// Ask user to register.
const exampleDone = () => {
  alert('buy me a coffee');
};
</script>

<style scoped>
.question {
  max-width: 1000px;
}

.choice-container {
  display: flex;
  align-items: center;
}

.desc-container {
  display: flex;
  align-items: center;
}

.highlight {
  border: 2px solid blue;
  border-radius: 5px;
  padding: 4px;
}

.question-id {
  cursor: pointer;
}
</style>

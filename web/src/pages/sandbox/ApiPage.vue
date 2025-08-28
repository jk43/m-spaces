<template>
  <q-page class="q-pa-md">
    <div class="text-h4 q-mb-md">API Test Page</div>

    <!-- 3x3x3 버튼 그리드 -->
    <div class="button-grid">
      <!-- 첫 번째 행 -->
      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="primary"
          label="rbac-manager"
          @click="handleButton111"
          class="col"
        />
        <q-btn
          color="secondary"
          label="Kafka history"
          @click="handleButton112"
          class="col"
        />
        <q-btn
          color="accent"
          label="Button 1-1-3"
          @click="handleButton113"
          class="col"
        />
      </div>

      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="positive"
          label="Button 1-2-1"
          @click="handleButton121"
          class="col"
        />
        <q-btn
          color="negative"
          label="Button 1-2-2"
          @click="handleButton122"
          class="col"
        />
        <q-btn
          color="warning"
          label="Button 1-2-3"
          @click="handleButton123"
          class="col"
        />
      </div>

      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="info"
          label="Button 1-3-1"
          @click="handleButton131"
          class="col"
        />
        <q-btn
          color="purple"
          label="Button 1-3-2"
          @click="handleButton132"
          class="col"
        />
        <q-btn
          color="orange"
          label="Button 1-3-3"
          @click="handleButton133"
          class="col"
        />
      </div>

      <!-- 두 번째 행 -->
      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="teal"
          label="Button 2-1-1"
          @click="handleButton211"
          class="col"
        />
        <q-btn
          color="pink"
          label="Button 2-1-2"
          @click="handleButton212"
          class="col"
        />
        <q-btn
          color="indigo"
          label="Button 2-1-3"
          @click="handleButton213"
          class="col"
        />
      </div>

      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="deep-orange"
          label="Button 2-2-1"
          @click="handleButton221"
          class="col"
        />
        <q-btn
          color="light-blue"
          label="Button 2-2-2"
          @click="handleButton222"
          class="col"
        />
        <q-btn
          color="lime"
          label="Button 2-2-3"
          @click="handleButton223"
          class="col"
        />
      </div>

      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="amber"
          label="Button 2-3-1"
          @click="handleButton231"
          class="col"
        />
        <q-btn
          color="cyan"
          label="Button 2-3-2"
          @click="handleButton232"
          class="col"
        />
        <q-btn
          color="deep-purple"
          label="Button 2-3-3"
          @click="handleButton233"
          class="col"
        />
      </div>

      <!-- 세 번째 행 -->
      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="brown"
          label="Button 3-1-1"
          @click="handleButton311"
          class="col"
        />
        <q-btn
          color="blue-grey"
          label="Button 3-1-2"
          @click="handleButton312"
          class="col"
        />
        <q-btn
          color="red"
          label="Button 3-1-3"
          @click="handleButton313"
          class="col"
        />
      </div>

      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="green"
          label="Button 3-2-1"
          @click="handleButton321"
          class="col"
        />
        <q-btn
          color="yellow"
          label="Button 3-2-2"
          @click="handleButton322"
          class="col"
        />
        <q-btn
          color="blue"
          label="Button 3-2-3"
          @click="handleButton323"
          class="col"
        />
      </div>

      <div class="row q-gutter-md q-mb-md">
        <q-btn
          color="grey"
          label="Button 3-3-1"
          @click="handleButton331"
          class="col"
        />
        <q-btn
          color="dark"
          label="Button 3-3-2"
          @click="handleButton332"
          class="col"
        />
        <q-btn
          color="white"
          text-color="black"
          label="Button 3-3-3"
          @click="handleButton333"
          class="col"
        />
      </div>
    </div>

    <!-- 결과 표시 영역 -->
    <q-card class="q-mt-lg">
      <q-card-section>
        <div class="text-h6">Button Click Result</div>
        <div class="text-body1">{{ lastClicked }}</div>
      </q-card-section>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref } from 'vue';
import { useQuasar } from 'quasar';
import { api } from 'src/boot/axios';

const $q = useQuasar();
const lastClicked = ref('No button has been clicked yet.');

// 첫 번째 행 함수들
const handleButton111 = () => {
  api()
    .get('/sandbox/rbac-manager')
    .then((res) => {
      console.log(res);
    });
};

const handleButton112 = () => {
  api()
    .get('/sandbox/kafka-history')
    .then((res) => {
      console.log(res);
    });
};

const handleButton113 = () => {
  lastClicked.value = 'Button 1-1-3 was clicked!';
  $q.notify({
    type: 'warning',
    message: 'Button 1-1-3 function executed',
    position: 'top',
  });
};

const handleButton121 = () => {
  lastClicked.value = 'Button 1-2-1 was clicked!';
  $q.notify({
    type: 'positive',
    message: 'Button 1-2-1 function executed',
    position: 'top',
  });
};

const handleButton122 = () => {
  lastClicked.value = 'Button 1-2-2 was clicked!';
  $q.notify({
    type: 'negative',
    message: 'Button 1-2-2 function executed',
    position: 'top',
  });
};

const handleButton123 = () => {
  lastClicked.value = 'Button 1-2-3 was clicked!';
  $q.notify({
    type: 'warning',
    message: 'Button 1-2-3 function executed',
    position: 'top',
  });
};

const handleButton131 = () => {
  lastClicked.value = 'Button 1-3-1 was clicked!';
  $q.notify({
    type: 'info',
    message: 'Button 1-3-1 function executed',
    position: 'top',
  });
};

const handleButton132 = () => {
  lastClicked.value = '버튼 1-3-2이 클릭되었습니다!';
  $q.notify({
    type: 'positive',
    message: '버튼 1-3-2 기능 실행',
    position: 'top',
  });
};

const handleButton133 = () => {
  lastClicked.value = '버튼 1-3-3이 클릭되었습니다!';
  $q.notify({
    type: 'warning',
    message: '버튼 1-3-3 기능 실행',
    position: 'top',
  });
};

// 두 번째 행 함수들
const handleButton211 = () => {
  lastClicked.value = '버튼 2-1-1이 클릭되었습니다!';
  $q.notify({
    type: 'positive',
    message: '버튼 2-1-1 기능 실행',
    position: 'top',
  });
};

const handleButton212 = () => {
  lastClicked.value = '버튼 2-1-2이 클릭되었습니다!';
  $q.notify({
    type: 'info',
    message: '버튼 2-1-2 기능 실행',
    position: 'top',
  });
};

const handleButton213 = () => {
  lastClicked.value = '버튼 2-1-3이 클릭되었습니다!';
  $q.notify({
    type: 'warning',
    message: '버튼 2-1-3 기능 실행',
    position: 'top',
  });
};

const handleButton221 = () => {
  lastClicked.value = '버튼 2-2-1이 클릭되었습니다!';
  $q.notify({
    type: 'positive',
    message: '버튼 2-2-1 기능 실행',
    position: 'top',
  });
};

const handleButton222 = () => {
  lastClicked.value = '버튼 2-2-2이 클릭되었습니다!';
  $q.notify({
    type: 'info',
    message: '버튼 2-2-2 기능 실행',
    position: 'top',
  });
};

const handleButton223 = () => {
  lastClicked.value = '버튼 2-2-3이 클릭되었습니다!';
  $q.notify({
    type: 'warning',
    message: '버튼 2-2-3 기능 실행',
    position: 'top',
  });
};

const handleButton231 = () => {
  lastClicked.value = '버튼 2-3-1이 클릭되었습니다!';
  $q.notify({
    type: 'positive',
    message: '버튼 2-3-1 기능 실행',
    position: 'top',
  });
};

const handleButton232 = () => {
  lastClicked.value = '버튼 2-3-2이 클릭되었습니다!';
  $q.notify({
    type: 'info',
    message: '버튼 2-3-2 기능 실행',
    position: 'top',
  });
};

const handleButton233 = () => {
  lastClicked.value = '버튼 2-3-3이 클릭되었습니다!';
  $q.notify({
    type: 'warning',
    message: '버튼 2-3-3 기능 실행',
    position: 'top',
  });
};

// 세 번째 행 함수들
const handleButton311 = () => {
  lastClicked.value = '버튼 3-1-1이 클릭되었습니다!';
  $q.notify({
    type: 'positive',
    message: '버튼 3-1-1 기능 실행',
    position: 'top',
  });
};

const handleButton312 = () => {
  lastClicked.value = '버튼 3-1-2이 클릭되었습니다!';
  $q.notify({
    type: 'info',
    message: '버튼 3-1-2 기능 실행',
    position: 'top',
  });
};

const handleButton313 = () => {
  lastClicked.value = '버튼 3-1-3이 클릭되었습니다!';
  $q.notify({
    type: 'warning',
    message: '버튼 3-1-3 기능 실행',
    position: 'top',
  });
};

const handleButton321 = () => {
  lastClicked.value = '버튼 3-2-1이 클릭되었습니다!';
  $q.notify({
    type: 'positive',
    message: '버튼 3-2-1 기능 실행',
    position: 'top',
  });
};

const handleButton322 = () => {
  lastClicked.value = '버튼 3-2-2이 클릭되었습니다!';
  $q.notify({
    type: 'info',
    message: '버튼 3-2-2 기능 실행',
    position: 'top',
  });
};

const handleButton323 = () => {
  lastClicked.value = '버튼 3-2-3이 클릭되었습니다!';
  $q.notify({
    type: 'warning',
    message: '버튼 3-2-3 기능 실행',
    position: 'top',
  });
};

const handleButton331 = () => {
  lastClicked.value = '버튼 3-3-1이 클릭되었습니다!';
  $q.notify({
    type: 'positive',
    message: '버튼 3-3-1 기능 실행',
    position: 'top',
  });
};

const handleButton332 = () => {
  lastClicked.value = '버튼 3-3-2이 클릭되었습니다!';
  $q.notify({
    type: 'info',
    message: '버튼 3-3-2 기능 실행',
    position: 'top',
  });
};

const handleButton333 = () => {
  lastClicked.value = '버튼 3-3-3이 클릭되었습니다!';
  $q.notify({
    type: 'warning',
    message: '버튼 3-3-3 기능 실행',
    position: 'top',
  });
};
</script>

<style scoped>
.button-grid {
  max-width: 800px;
  margin: 0 auto;
}

.q-btn {
  min-height: 60px;
  font-size: 14px;
}
</style>

<template>
  <div>
    <vue-avatar
      :width="400"
      :height="400"
      :rotation="Number(rotation)"
      :borderRadius="Number(borderRadius)"
      :scale="Number(scale)"
      :image="profileImage"
      ref="vueavatar"
      @vue-avatar-editor:image-ready="onImageReady"
    >
    </vue-avatar>
    <br />
    <label>
      Zoom : {{ scale }}x
      <br />
      <q-slider :min="1" :max="3" :step="0.02" v-model="scale" />
    </label>
    <br />
    <label>
      Rotation : {{ rotation }}°
      <br />
      <q-slider :min="0" :max="360" :step="1" v-model="rotation" />
    </label>
    <br />
    <button v-on:click="saveClicked">Get image</button>
    <br />
    <img ref="image" />
  </div>
  <div>
    <q-file
      v-model="files"
      label="Pick files"
      counter
      filled
      use-chips
      multiple
      append
      @input="uploadFiles"
    />
    <button v-on:click="uploadFiles">upload</button>
    {{ files }}
  </div>

  <div>
    <q-form @submit="onSubmit" @reset="onReset" class="q-gutter-md">
      <q-input
        filled
        v-model="form.name"
        label="Name"
        lazy-rules
        :rules="[(val) => !!val || 'Field is required']"
      />
      <q-input
        filled
        v-model="form.email"
        label="Email"
        lazy-rules
        :rules="[(val) => !!val || 'Field is required']"
      />
      <q-file
        v-model="form.picture"
        label="Picture"
        filled
        lazy-rules
        :rules="[(val) => !!val || 'Field is required']"
      />
      <q-file
        v-model="form.video"
        label="Video"
        filled
        lazy-rules
        :rules="[(val) => !!val || 'Field is required']"
      />
      <div>
        <q-btn label="Submit" type="submit" color="primary" />
        <q-btn
          label="Reset"
          type="reset"
          color="primary"
          flat
          class="q-ml-sm"
        />
      </div>
    </q-form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { VueAvatar } from 'vue-avatar-editor-improved';
import { api } from 'src/boot/axios';
import { Payload, FormDataRequst } from 'src/types';
import { useUserStore } from 'src/stores/user-store';
import { refreshToken } from 'src/service/user/auth';

const rotation = ref<number>(0);
const scale = ref<number>(1);
const borderRadius = ref<number>(200);

const vueavatar = ref<VueAvatar | null>(null);
const image = ref<HTMLImageElement | null>();

const files = ref<FileList | null>();

const onImageReady = () => {
  scale.value = 1;
  rotation.value = 0;
};

const form = ref({
  name: '',
  email: '',
  picture: null,
  video: null,
});

const profileImage = useUserStore().getProfileImage;

const onSubmit = async () => {
  const formData = new FormData();
  formData.append('name', form.value.name);
  formData.append('email', form.value.email);
  formData.append('picture', form.value.picture);
  formData.append('video', form.value.video);

  try {
    const response = await api().post('/files', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    console.log(response.data);
  } catch (err) {
    console.error(err);
  }
};

const saveClicked = async () => {
  const newImage = vueavatar.value.getImageScaled();
  if (newImage && image.value) {
    console.log('iffffff: ');
    image.value.src = newImage.toDataURL();

    // Convert the image data URL to a Blob
    const response = await fetch(image.value.src);
    const blob = await response.blob();
    const newFile = new File([blob], 'newFileName.jpg', {
      type: blob.type,
    });

    // Create a FormData object and append the image blob
    const formData = new FormData();
    formData.append('file', newFile);
    formData.append('service', 'user_service');
    formData.append('serviceCtx', 'profile-image');

    // Send the image data to the server
    await api().post('/files', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  }
  refreshToken(false, null);
};

const uploadFiles = async () => {
  if (files.value) {
    const formData = new FormData();
    const payload: Payload<FormDataRequst> = {
      data: {
        file: null,
        service: 'user_service',
        serviceCtx: 'profile-image',
      },
    };
    for (let i = 0; i < files.value.length; i++) {
      formData.append('file', files.value[i]);
      formData.append('service', 'user_service');
      formData.append('servieCtx', 'profile-image');
    }

    // Send the image data to the server
    api().post('/files', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  }
};
</script>

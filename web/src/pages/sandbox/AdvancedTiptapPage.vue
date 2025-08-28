<template>
  <q-page class="q-pa-md">
    <div class="row q-col-gutter-md">
      <div class="col-12">
        <q-card>
          <q-card-section>
            <div class="text-h6">Advanced Tiptap Editor Example</div>
            <div class="text-caption text-grey-6">
              Advanced rich text editor with various plugins included.
            </div>
          </q-card-section>

          <q-card-section>
            <div class="row q-col-gutter-md">
              <div class="col-12">
                <div class="text-subtitle2 q-mb-sm">Advanced Editor</div>
                <AdvancedTiptapEditor
                  v-model="content"
                  placeholder="Enter your content here..."
                  height="500px"
                  @change="handleChange"
                />
              </div>
            </div>
          </q-card-section>

          <q-card-section>
            <div class="row q-col-gutter-md">
              <div class="col-12 col-md-6">
                <div class="text-subtitle2 q-mb-sm">HTML Output</div>
                <q-input
                  v-model="content"
                  type="textarea"
                  outlined
                  readonly
                  rows="10"
                  class="q-mt-sm"
                />
              </div>

              <div class="col-12 col-md-6">
                <div class="text-subtitle2 q-mb-sm">Preview</div>
                <div class="preview-content q-pa-md" v-html="content"></div>
              </div>
            </div>
          </q-card-section>

          <q-card-section>
            <div class="row q-col-gutter-sm">
              <div class="col-auto">
                <q-btn
                  color="primary"
                  label="Add Sample Document"
                  @click="addSampleDocument"
                />
              </div>
              <div class="col-auto">
                <q-btn
                  color="secondary"
                  label="Clear Content"
                  @click="clearContent"
                />
              </div>
              <div class="col-auto">
                <q-btn color="accent" label="Copy HTML" @click="copyHTML" />
              </div>
            </div>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useQuasar } from 'quasar';
import AdvancedTiptapEditor from 'src/components/AdvancedTiptapEditor.vue';

const $q = useQuasar();

const content = ref('');

const handleChange = (value: string) => {
  console.log('에디터 내용 변경:', value);
};

const clearContent = () => {
  content.value = '';
  $q.notify({
    type: 'positive',
    message: 'Content has been cleared.',
  });
};

const addSampleDocument = () => {
  content.value = `
    <h1>Advanced Tiptap Editor Sample Document</h1>

    <h2>Text Formatting</h2>
    <p>This is a sample that includes <strong>bold text</strong> and <em>italic text</em>.</p>
    <p><u>Underlined text</u> and <s>strikethrough text</s> can also be used.</p>
    <p><code>Inline code</code> is also supported.</p>

    <h2>Colors and Highlights</h2>
    <p>You can change text color and background color.</p>

    <h2>Lists</h2>
    <ul>
      <li>Bullet list item 1</li>
      <li>Bullet list item 2</li>
      <li>Bullet list item 3</li>
    </ul>

    <ol>
      <li>Numbered list item 1</li>
      <li>Numbered list item 2</li>
      <li>Numbered list item 3</li>
    </ol>

    <h2>Checklist</h2>
    <ul data-type="taskList">
      <li data-type="taskItem" data-checked="true">
        <label>
          <input type="checkbox" checked="checked">
        </label>
        <div><p>Completed task</p></div>
      </li>
      <li data-type="taskItem" data-checked="false">
        <label>
          <input type="checkbox">
        </label>
        <div><p>Incomplete task</p></div>
      </li>
    </ul>

    <h2>Blockquote</h2>
    <blockquote>
      This is a blockquote. It is used to emphasize important content.
    </blockquote>

    <h2>Code Block</h2>
    <pre><code class="language-javascript">function hello() {
  console.log('Hello, Tiptap!');
}

hello();</code></pre>

    <h2>Table</h2>
    <table>
      <thead>
        <tr>
          <th>Header 1</th>
          <th>Header 2</th>
          <th>Header 3</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>Content 1</td>
          <td>Content 2</td>
          <td>Content 3</td>
        </tr>
        <tr>
          <td>Content 4</td>
          <td>Content 5</td>
          <td>Content 6</td>
        </tr>
      </tbody>
    </table>

    <h2>Links and Images</h2>
    <p>Visit the <a href="https://tiptap.dev">Tiptap official site</a>.</p>
    <p>You can also add images:</p>
    <p><img src="https://via.placeholder.com/300x200" alt="Sample image"></p>

    <h2>Alignment</h2>
    <p style="text-align: left;">This is left-aligned text.</p>
    <p style="text-align: center;">This is center-aligned text.</p>
    <p style="text-align: right;">This is right-aligned text.</p>
    <p style="text-align: justify;">This is justified text. This text is aligned to both ends for better appearance.</p>
  `;
  $q.notify({
    type: 'positive',
    message: 'Sample document has been added.',
  });
};

const copyHTML = () => {
  navigator.clipboard
    .writeText(content.value)
    .then(() => {
      $q.notify({
        type: 'positive',
        message: 'HTML has been copied to clipboard.',
      });
    })
    .catch(() => {
      $q.notify({
        type: 'negative',
        message: 'Failed to copy to clipboard.',
      });
    });
};
</script>

<style scoped></style>

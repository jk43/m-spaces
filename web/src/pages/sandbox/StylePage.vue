<template>
  <div
    class="q-pa-lg q-mt-lg bg-blue-grey flex items-center justify-center q-gutter-md"
  >
    <span class="flex-center text-blue-grey-1 shadow-up-15 bg-primary">
      StylePage
    </span>
    <span class="flex-center text-blue-grey-1 shadow-up-15 bg-primary">
      StylePage
    </span>
    <span class="flex-center text-blue-grey-1 shadow-up-15 bg-primary">
      StylePage
    </span>
  </div>
  <div class="example-row-equal-width">
    <div class="row">
      <div class="col">.col</div>
      <div class="col">.col</div>
    </div>

    <div class="row">
      <div class="col">.col</div>
      <div class="col">.col</div>
      <div class="col">.col</div>
    </div>
  </div>
  <p>Flex Row</p>
  <div class="my-row example-row-equal-width row q-gutter-lg">
    <div class="col text-center">
      <q-btn label="hello" class="text-caption text-left"></q-btn>
    </div>
    <div class="col text-right">
      <q-btn label="hello" class="text-primary text-capitalize"></q-btn>
    </div>
    <div class="col"><q-btn label="hello"></q-btn></div>
    <div class="col"><q-btn label="hello"></q-btn></div>
  </div>
</template>

<style lang="sass">
.example-row-equal-width
  .row > div
    background: rgba(#999,.15)
    border: 1px solid rgba(#999,.2)
  .row + .row
    margin-top: 1rem
.my-row1
  .col
    background: rgba(#999,.15)
    border: 1px solid rgba(#999,.2)
    padding: 10px 15px
</style>

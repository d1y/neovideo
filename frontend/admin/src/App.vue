<script setup lang="ts">
import { onMounted, ref } from 'vue';

const msg = ref("hello world")

onMounted(()=> {
  msg.value = "hello world 1"
})

function d() {
  debugger
}
</script>

<template>
  <div @click="d"> {{ msg }} </div>
</template>
<script setup lang="ts">
import { ref } from "vue";

const { text, defaultOpen = false } = defineProps<{
    text: string;
    defaultOpen?: boolean;
}>();
const show = ref(defaultOpen);
</script>

<template>
    <button class="btn-primary btn-rounded" @click="show = !show">
        {{ text }}
        <i
            class="bi"
            :class="{
                'bi-caret-down-fill': !show,
                'bi-caret-up-fill': show,
            }"
        ></i>
    </button>
    <slot v-if="show"></slot>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
</style>

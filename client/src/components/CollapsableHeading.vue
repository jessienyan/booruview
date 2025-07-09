<script setup lang="ts">
import { ref } from "vue";

const {
    el,
    heading,
    defaultOpen = false,
} = defineProps<{ el: string; heading: string; defaultOpen?: boolean }>();
const show = ref(defaultOpen);
</script>

<template>
    <component :is="el" @click="show = !show"
        ><button class="btn-toggle">
            {{ heading }}
            <i
                class="bi"
                :class="{
                    'bi-caret-down-fill': !show,
                    'bi-caret-up-fill': show,
                }"
            ></i></button
    ></component>
    <slot v-if="show"></slot>
</template>

<style lang="scss" scoped>
.btn-toggle {
    background: none;
    border: none;
    color: inherit;
    font-size: inherit;
    font-weight: inherit;
    cursor: pointer;
    padding: 0;
}
</style>

<script setup lang="ts">
import { ref } from "vue";

const {
	text,
	button = true,
	defaultOpen = false,
} = defineProps<{
	text: string;
	button?: boolean;
	defaultOpen?: boolean;
}>();
const show = ref(defaultOpen);
</script>

<template>
    <button v-if="button" class="btn-primary btn-rounded" @click="show = !show">
        {{ text }}
        <i
            class="bi"
            :class="{
                'bi-caret-down-fill': !show,
                'bi-caret-up-fill': show,
            }"
        ></i>
    </button>
    <button v-else class="btn-unstyled" @click="show = !show">

        <i
            class="bi"
            :class="{
                'bi-caret-right-fill': !show,
                'bi-caret-down-fill': show,
            }"
        ></i>
    
        {{ text }}
    </button>
    <slot v-if="show"></slot>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-unstyled {
    display: block;
    border: none;
    background: none;
    padding: none;
    color: inherit;
    cursor: pointer;
    margin: 1em 0;
}
</style>

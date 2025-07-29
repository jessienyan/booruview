<script setup lang="ts">
import type store from "@/store";

type ToastKind = typeof store.toast.type;
const { kind } = defineProps<{ kind: ToastKind }>();
defineEmits(["dismiss"]);
</script>

<template>
    <div
        class="toast"
        :class="{ [`toast-${kind}`]: true }"
        @click="$emit('dismiss')"
    >
        <slot></slot>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

.toast {
    position: absolute;
    bottom: 80px;
    left: 50%;
    transform: translateX(-50%);
    max-width: 90vw;
    width: max-content;
    z-index: 999;
    border-radius: 0.8rem;
    padding: 1rem 1.2rem;
    cursor: pointer;
    filter: drop-shadow(0 0 10px black);
    color: #fff;
}

.toast-info {
    background-color: $color-primary-lighter;
}

.toast-error {
    background-color: #700;
}
</style>

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
    bottom: 40px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 999;
    border-radius: 0.8rem;
    padding: 0.8rem 1rem;
    cursor: pointer;
    filter: drop-shadow(0 0 10px black);
}

.toast-info {
    background-color: $color-primary-lighter;
}

.toast-error {
    background-color: #700;
}
</style>

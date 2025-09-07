<script setup lang="ts">
import type store from "@/store";
import { onMounted, ref } from "vue";

type ToastKind = typeof store.toast.type;
const { kind } = defineProps<{ kind: ToastKind }>();
const emit = defineEmits(["dismiss"]);
const dismissing = ref(false);
const appearing = ref(true);

function onDismiss() {
    setTimeout(() => emit("dismiss"), 400);
    dismissing.value = true;
}

onMounted(() => {
    setTimeout(() => (appearing.value = false), 400);
});
</script>

<template>
    <div
        class="toast"
        :class="{ [`toast-${kind}`]: true, appearing, dismissing }"
        @click="onDismiss"
    >
        <slot></slot>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

@keyframes slide-up {
    0% {
        bottom: 55px;
        opacity: 0;
    }
    100% {
        bottom: 80px;
        opacity: 1;
    }
}

.toast {
    position: absolute;
    bottom: 80px;
    left: 50%;
    transform: translateX(-50%);
    max-width: 90vw;
    width: max-content;
    z-index: 999;
    border-radius: 0.8rem;
    padding: 1.2rem 1.6rem;
    font-size: 1.1rem;
    cursor: pointer;
    filter: drop-shadow(0 0 10px black);
    color: #fff;
    text-shadow: 0 0 3px #000a;
}

.appearing {
    animation: 300ms ease-out slide-up;
}

.dismissing {
    animation: 180ms ease-out reverse slide-up;
    opacity: 0;
}

.toast-info {
    background-color: $color-primary-lighter;
}

.toast-error {
    background-color: #700;
}
</style>

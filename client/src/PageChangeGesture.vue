<script setup lang="ts">
import VanillaSwipe, { type EventData } from "vanilla-swipe";
import { onMounted, onUnmounted, ref } from "vue";

const props = defineProps<{ scrollContainer: HTMLElement }>();
const swipeDirection = ref<"LEFT" | "RIGHT" | null>(null);
const minDistanceForSwipe = 50; // pixels
const maxAngleForSwipe = 30; // degrees

let swipe: VanillaSwipe;

function isHorizontalSwipe({ absY, absX }: EventData): boolean {
    const angle = Math.atan2(absY, absX) * (180 / Math.PI);
    return angle < maxAngleForSwipe;
}

onMounted(() => {
    swipe = new VanillaSwipe({
        element: props.scrollContainer,
        delta: minDistanceForSwipe,
        directionDelta: 30,

        // Swap the swipe direction (flick left means right)
        rotationAngle: 180,

        onSwipeStart(event, data) {
            if (isHorizontalSwipe(data)) {
                swipeDirection.value = data.directionX as any;
            } else {
                swipeDirection.value = null;
            }
        },

        onSwiping(event, data) {
            if (swipeDirection.value === null) {
                return;
            }

            // Cancel the swipe if the user changed direction
            if (swipeDirection.value !== data.directionX) {
                swipeDirection.value = null;
            }
        },

        onSwiped(event, data) {
            if (swipeDirection.value === null) {
                return;
            }
        },
    });
    swipe.init();
});

onUnmounted(() => swipe.destroy());
</script>

<template>
    <div class="swipe swipe-left" v-if="swipeDirection === 'LEFT'">
        <i class="bi bi-arrow-left"></i>
    </div>
    <div class="swipe swipe-right" v-if="swipeDirection === 'RIGHT'">
        <i class="bi bi-arrow-right"></i>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

.swipe {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    z-index: 1;

    padding: 20px;
    margin: 10px;

    background-color: $color-primary;
    color: $color-primary-light;
    font-size: 36px;
    border-radius: 100%;
    filter: drop-shadow(0 0 10px black);
    opacity: 0.8;
}

@keyframes slide-from-left {
    0% {
        left: -100px;
        opacity: 0;
    }

    100% {
        left: 0;
    }
}

@keyframes slide-from-right {
    0% {
        right: -100px;
        opacity: 0;
    }

    100% {
        right: 0;
    }
}

.swipe-left {
    animation: slide-from-left 400ms cubic-bezier(0, 1, 0.25, 1) 1;
    left: 0;
}

.swipe-right {
    animation: slide-from-right 400ms cubic-bezier(0, 1, 0.25, 1) 1;
    right: 0;
}
</style>

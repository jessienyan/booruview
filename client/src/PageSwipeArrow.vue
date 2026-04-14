<script setup lang="ts">
import VanillaSwipe, { type EventData } from "vanilla-swipe";
import { onMounted, onUnmounted, ref, watch } from "vue";
import store from "./store";

const props = defineProps<{
    scrollContainer: HTMLElement;
    currentPage: number;
    maxPage: number;
}>();

const emit = defineEmits(["prev", "next"]);

const swipeDirection = ref<"LEFT" | "RIGHT" | null>(null);
const minDistanceForSwipe = 50; // pixels
const maxAngleForSwipe = 30; // degrees

let swipe: VanillaSwipe;

function isHorizontalSwipe({ absY, absX }: EventData): boolean {
    const angle = Math.atan2(absY, absX) * (180 / Math.PI);
    return angle < maxAngleForSwipe;
}

watch(swipeDirection, () => {
    store.userIsSwipingToChangePage = swipeDirection.value !== null;
});

onMounted(() => {
    swipe = new VanillaSwipe({
        element: props.scrollContainer,
        delta: minDistanceForSwipe,
        directionDelta: 10,

        // Swap the swipe direction (flick left means right)
        rotationAngle: 180,

        onSwipeStart(_, data) {
            if (store.fullscreenPost !== null || !isHorizontalSwipe(data)) {
                swipeDirection.value = null;
                return;
            }

            switch (data.directionX) {
                case "LEFT":
                    if (props.currentPage <= 1) {
                        return;
                    }

                    swipeDirection.value = "LEFT";
                    break;

                case "RIGHT":
                    if (props.currentPage === props.maxPage) {
                        return;
                    }

                    swipeDirection.value = "RIGHT";
                    break;
            }
        },

        onSwiping(_, data) {
            if (swipeDirection.value === null) {
                return;
            }

            // Cancel the swipe if the user changed direction
            if (swipeDirection.value !== data.directionX) {
                swipeDirection.value = null;
                return;
            }
        },

        onSwiped() {
            switch (swipeDirection.value) {
                case "LEFT":
                    emit("prev");
                    swipeDirection.value = null;
                    break;
                case "RIGHT":
                    emit("next");
                    swipeDirection.value = null;
                    break;
            }
        },
    });
    swipe.init();
});

onUnmounted(() => swipe.destroy());
</script>

<template>
    <div class="page-swipe-container">
        <div class="swipe swipe-left" v-if="swipeDirection === 'LEFT'">
            <i class="bi bi-arrow-left"></i>
        </div>
        <div class="swipe swipe-right" v-if="swipeDirection === 'RIGHT'">
            <i class="bi bi-arrow-right"></i>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

.page-swipe-container {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 1;
    pointer-events: none;
    user-select: none;
}

.swipe {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    padding: 15px;
    margin: 0 20px;
    background-color: $color-primary-light;
    color: $color-primary;
    font-size: 36px;
    line-height: 0;
    border-radius: 100%;
    filter: drop-shadow(0 0 10px black);
    opacity: 0.9;
}

@keyframes slide-from-left {
    0% {
        left: -100px;
    }

    100% {
        left: 0;
    }
}

@keyframes slide-from-right {
    0% {
        right: -100px;
    }

    100% {
        right: 0;
    }
}

.swipe-left {
    animation: slide-from-left 200ms cubic-bezier(0, 1, 0.25, 1) 1;
    left: 0;
}

.swipe-right {
    animation: slide-from-right 200ms cubic-bezier(0, 1, 0.25, 1) 1;
    right: 0;
}
</style>

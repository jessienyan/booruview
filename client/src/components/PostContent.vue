<script setup lang="ts">
import { useIsVideo } from "@/composable";
import store from "@/store";
import { computed, onUnmounted, ref, useTemplateRef, watch } from "vue";

const {
    cropped,
    maxHeight,
    renderHeight,
    post,
    scrollContainer,
    beingDragged,
} = defineProps<{
    cropped: boolean;
    maxHeight: number;
    renderHeight: number;
    post: Post;
    scrollContainer: HTMLElement;
    beingDragged?: boolean;
}>();

const isVideo = useIsVideo(() => post);
const favorited = computed(
    () => store.settings.favorites.findIndex((p) => p.id === post.id) !== -1,
);

const content = computed<{ url: string; width: number; height: number }>(() => {
    if (isVideo.value) {
        return {
            url: post.thumbnail_url,
            width: post.thumbnail_width,
            height: post.thumbnail_height,
        };
    }

    if (post.lowres_url.length > 0) {
        return {
            url: post.lowres_url,
            width: post.lowres_width,
            height: post.lowres_height,
        };
    }

    return {
        url: post.image_url,
        width: post.width,
        height: post.height,
    };
});

const showImage = ref(false);
const containerRef = useTemplateRef("container");
const scrollObserver = new IntersectionObserver(onIntersectionChange, {
    root: scrollContainer,

    // Preload images when they are within this distance from the bottom of the viewport
    rootMargin: "600px 0px 600px 0px",
});

function onIntersectionChange(entries: IntersectionObserverEntry[]) {
    const e = entries[0];

    // Once: render the image when it comes into view and cleanup the observer
    if (e.isIntersecting) {
        showImage.value = true;
        scrollObserver.disconnect();
    }
}

// Reset showImage state if post changes. Allows component to be reused
watch(
    () => [post.id, containerRef.value],
    () => {
        if (containerRef.value === null) {
            return;
        }

        showImage.value = false;
        scrollObserver.observe(containerRef.value!);
    },
);

onUnmounted(() => {
    scrollObserver.disconnect();
});
</script>

<template>
    <div
        class="post"
        :style="{ maxHeight: maxHeight + 'px', height: renderHeight + 'px' }"
        @click="store.fullscreenPost = post"
        ref="container"
        :draggable="true"
    >
        <img
            v-if="showImage"
            class="content"
            referrerpolicy="no-referrer"
            :src="content.url"
            :width="content.width"
            :height="content.height"
        />

        <div class="drag-container" v-if="beingDragged">
            <div class="drag-icon-container">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="50"
                    height="50"
                    fill="currentColor"
                    class="bi bi-arrows-move drag-icon"
                    viewBox="0 0 16 16"
                >
                    <path
                        fill-rule="evenodd"
                        d="M7.646.146a.5.5 0 0 1 .708 0l2 2a.5.5 0 0 1-.708.708L8.5 1.707V5.5a.5.5 0 0 1-1 0V1.707L6.354 2.854a.5.5 0 1 1-.708-.708zM8 10a.5.5 0 0 1 .5.5v3.793l1.146-1.147a.5.5 0 0 1 .708.708l-2 2a.5.5 0 0 1-.708 0l-2-2a.5.5 0 0 1 .708-.708L7.5 14.293V10.5A.5.5 0 0 1 8 10M.146 8.354a.5.5 0 0 1 0-.708l2-2a.5.5 0 1 1 .708.708L1.707 7.5H5.5a.5.5 0 0 1 0 1H1.707l1.147 1.146a.5.5 0 0 1-.708.708zM10 8a.5.5 0 0 1 .5-.5h3.793l-1.147-1.146a.5.5 0 0 1 .708-.708l2 2a.5.5 0 0 1 0 .708l-2 2a.5.5 0 0 1-.708-.708L14.293 8.5H10.5A.5.5 0 0 1 10 8"
                    />
                </svg>
            </div>
        </div>

        <div class="icons" v-if="favorited || cropped">
            <i
                v-if="favorited"
                class="bi bi-heart-fill fav-icon"
                title="favorited"
            ></i>
            <i
                v-if="cropped"
                class="bi bi-crop crop-icon"
                title="post is cropped (too tall)"
            ></i>
        </div>
        <i v-if="isVideo" class="bi bi-play-circle play-icon"></i>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

.post {
    break-inside: avoid;
    font-size: 0;
    position: relative;
    cursor: pointer;
}

.content {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.icons {
    position: absolute;
    right: 0.8rem;
    top: 0.8rem;
    color: white;
    background-color: rgba(0, 0, 0, 0.75);
    opacity: 0.5;
    border-radius: 9999px;
    padding: 0.5em;
    font-size: 26px;
    display: flex;
    flex-direction: column;
    gap: 0.5em;

    .bi {
        width: 26px;
        height: 26px;
        text-shadow: 0 0 3px black;
    }

    .crop-icon {
        position: relative;
        scale: 0.9;
        top: -0.12em;
        left: -0.06em;
    }

    .fav-icon {
        color: pink;
    }
}

.play-icon {
    position: absolute;
    font-size: 60px;
    left: 50%;
    top: 50%;
    transform: translateX(-50%) translateY(-50%);
    filter: drop-shadow(0 0 3px black);
    color: #fff;
    opacity: 0.6;
}

.drag-container {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    border: 5px solid yellow;
}

.drag-icon-container {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translateX(-50%) translateY(-50%);
    padding: 15px;
    border-radius: 50%;
    background-color: rgba(0, 0, 0, 0.5);
}

.drag-icon {
    filter: drop-shadow(0 0 3px black);
    color: #eee;
}
</style>

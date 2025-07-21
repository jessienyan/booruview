<script setup lang="ts">
import { useIsVideo } from "@/composable";
import store from "@/store";
import { computed, onMounted, onUnmounted, ref, useTemplateRef } from "vue";

const { cropped, maxHeight, renderHeight, post, scrollContainer } =
    defineProps<{
        cropped: boolean;
        maxHeight: number;
        renderHeight: number;
        post: Post;
        scrollContainer: HTMLElement;
    }>();

const isVideo = useIsVideo(() => post);

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

onMounted(() => scrollObserver.observe(containerRef.value!));
onUnmounted(() => scrollObserver.disconnect());
</script>

<template>
    <div
        class="post"
        :style="{ maxHeight: maxHeight + 'px', height: renderHeight + 'px' }"
        @click="store.fullscreenPost = post"
        ref="container"
    >
        <img
            class="content"
            :src="content.url"
            :width="content.width"
            :height="content.height"
            v-if="showImage"
        />

        <span
            v-if="cropped"
            class="crop-icon"
            title="post is cropped (too tall)"
            ><i class="bi bi-crop"></i
        ></span>
        <i v-if="isVideo" class="bi bi-play-circle play-icon"></i>
    </div>
</template>

<style lang="scss" scoped>
.post {
    break-inside: avoid;
    font-size: 0;
    position: relative;
    cursor: pointer;
}

.content {
    /* placeholder color */
    background-color: #444;
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.crop-icon {
    font-size: 24px;
    position: absolute;
    right: 0.8rem;
    top: 0.8rem;
    background-color: black;
    color: white;
    opacity: 0.25;
    border-radius: 50%;
    padding: 0.6em;
    line-height: 0;

    .bi {
        display: block;
        line-height: 0;
        position: relative;
        left: -0.06em;
        top: -0.06em;
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
</style>

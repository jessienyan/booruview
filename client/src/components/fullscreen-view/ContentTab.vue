<script setup lang="ts">
import { computed, onMounted, onUnmounted, useTemplateRef } from "vue";
import {
    useGelbooruImageURL,
    useGelbooruVideoURL,
    useIsVideo,
    usePanZoom,
    useStationaryClick,
} from "@/composable";
import store from "@/store";

const imgRef = useTemplateRef("imgRef");
const { post } = defineProps<{ post: Post }>();
const htmlRoot = document.body.parentElement as HTMLElement;
const overscrollCssClass = "prevent-overscroll";
const isVideo = useIsVideo(() => post);
const emit = defineEmits(["prev", "next"]);

const content = computed(() => {
    const hasHighRes = post.image_url.length > 0;
    const hasLowRes = post.lowres_url.length > 0;

    // Avoid using high res images if the media proxy is enabled
    const useHighRes =
        !hasLowRes ||
        (!store.cdnHosts?.media_proxy &&
            hasHighRes &&
            store.settings.highResImages);

    if (useHighRes) {
        return {
            url: post.image_url,
            height: post.height,
            width: post.width,
        };
    }

    return {
        url: post.lowres_url,
        height: post.lowres_height,
        width: post.lowres_width,
    };
});

const imageURL = useGelbooruImageURL(() => content.value.url);
const videoURL = useGelbooruVideoURL(() => content.value.url);

usePanZoom({
    enable: () => store.settings.enablePanZoom && !isVideo.value,
    el: imgRef,
    key: () => post.id,
});

const goPrev = useStationaryClick(() => emit("prev"));
const goNext = useStationaryClick(() => emit("next"));

onMounted(() => {
    // Since the touch event isn't being blocked we need to prevent the user from
    // overscrolling the page (refresh by pulling down)
    htmlRoot.classList.add(overscrollCssClass);
});

onUnmounted(() => {
    htmlRoot.classList.remove(overscrollCssClass);
});
</script>

<template>
    <video
        v-if="isVideo"
        :poster="imageURL"
        :autoplay="store.settings.autoplayVideo"
        :muted="store.settings.muteVideo"
        :key="`video-${post.id}`"
        controls
        loop
    >
        <source
            :src="videoURL"
            type="video/mp4"
            v-if="videoURL.endsWith('.mp4')"
        />
        <source
            :src="videoURL"
            type="video/webm"
            v-if="videoURL.endsWith('.webm')"
        />
    </video>

    <template v-else>
        <template v-if="store.settings.enableClickImageToChange">
            <div
                class="hidden-prev-btn"
                @mousedown="goPrev.mouseDown"
                @mouseup="goPrev.mouseUp"
            ></div>
            <div
                class="hidden-next-btn"
                @mousedown="goNext.mouseDown"
                @mouseup="goNext.mouseUp"
            ></div>
        </template>

        <img
            v-if="store.settings.enablePanZoom"
            ref="imgRef"
            referrerpolicy="same-origin"
            :src="imageURL"
            :width="content.width"
            :height="content.height"
            :key="post.id"
        />
        <img
            v-else
            class="img-fit"
            ref="imgRef"
            referrerpolicy="same-origin"
            :src="imageURL"
        />
    </template>
</template>

<style scoped>
video {
    max-width: 100%;
    max-height: 100%;
}

.img-fit {
    max-width: 100%;
    max-height: 100%;
    transform: translateY(-50%);
    position: relative;
    top: 50%;
}

.hidden-prev-btn {
    z-index: 1;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;

    /* prev button takes up 33% of the screen */
    right: 67%;
}

.hidden-next-btn {
    z-index: 1;
    position: absolute;
    top: 0;
    bottom: 0;

    /* next button takes up 67% of the screen. this allows users to click either
       the right side of the viewport or the center
    */
    left: 33%;
    right: 0;
}
</style>

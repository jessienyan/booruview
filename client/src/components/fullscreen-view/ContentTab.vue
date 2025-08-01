<script setup lang="ts">
import store from "@/store";
import {
    computed,
    onActivated,
    onDeactivated,
    onMounted,
    onUnmounted,
    useTemplateRef,
    watch,
} from "vue";
import createPanZoom, { type PanZoom } from "panzoom";
import { useIsVideo } from "@/composable";

const imgRef = useTemplateRef("imgRef");
let pz: PanZoom | undefined;
const { post } = defineProps<{ post: Post }>();
const htmlRoot = document.body.parentElement as HTMLElement;
const overscrollCssClass = "prevent-overscroll";
const isVideo = useIsVideo(() => post);

function setupPanZoom() {
    pz?.dispose();

    // Don't use panzoom for videos
    if (isVideo.value) {
        return;
    }

    pz = createPanZoom(imgRef.value!, {
        autocenter: true,
        bounds: true,
        boundsPadding: 0.1,
        maxZoom: 4,
        minZoom: 0.05,
        onTouch() {
            // Don't block the touch event so the user can right click
            return false;
        },
    });
}

watch(() => post.id, setupPanZoom, { flush: "post" });

onMounted(() => {
    setupPanZoom();

    // Since the touch event isn't being blocked we need to prevent the user from
    // overscrolling the page (refresh by pulling down)
    htmlRoot.classList.add(overscrollCssClass);
});

onUnmounted(() => {
    pz?.dispose();
    htmlRoot.classList.remove(overscrollCssClass);
});

onDeactivated(() => pz?.pause());
onActivated(() => pz?.resume());

const content = computed(() => {
    const hasHighRes = post.image_url.length > 0;
    const hasLowRes = post.lowres_url.length > 0;

    if (!hasLowRes || (hasHighRes && store.settings.highResImages)) {
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
</script>

<template>
    <video
        v-if="isVideo"
        :poster="post.thumbnail_url || post.lowres_url"
        :autoplay="store.settings.autoplayVideo"
        :muted="store.settings.muteVideo"
        controls
        loop
    >
        <source
            :src="content.url"
            type="video/mp4"
            v-if="content.url.endsWith('.mp4')"
        />
        <source
            :src="content.url"
            type="video/webm"
            v-if="content.url.endsWith('.webm')"
        />
    </video>

    <!-- Using a key on the image prevents it from stretching when changing between posts -->
    <img
        v-else
        :src="content.url"
        :width="content.width"
        :height="content.height"
        :key="post.id"
        ref="imgRef"
    />
</template>

<style scoped>
video {
    max-width: 100%;
    max-height: 100%;
}
</style>

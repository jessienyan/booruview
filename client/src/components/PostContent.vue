<script setup lang="ts">
import { computed, onUnmounted, ref, useTemplateRef, watch } from "vue";
import { useGelbooruImageURL, useIsVideo } from "@/composable";
import store from "@/store";

const { cropped, maxHeight, renderHeight, post, scrollContainer, forceLoad } =
    defineProps<{
        cropped: boolean;
        maxHeight: number;
        renderHeight: number;
        post: Post;
        scrollContainer: HTMLElement;

        // When true, the post's content will be loaded regardless if it's on screen.
        // Can be used for preloading posts
        forceLoad?: boolean;
    }>();

const isVideo = useIsVideo(() => post);
const favPosts = store.favoritePosts();
const favorited = computed(
    () => favPosts.value.findIndex((p) => p.id === post.id) !== -1,
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

const imageURL = useGelbooruImageURL(() => content.value.url);

const showImage = ref(false);
const containerRef = useTemplateRef("container");

// Load images that are close to being visible
const preloadImageViewportDistance = "1200px";
const scrollObserver = new IntersectionObserver(onIntersectionChange, {
    root: scrollContainer,
    rootMargin: `${preloadImageViewportDistance} 0px ${preloadImageViewportDistance} 0px`,
});

watch(
    () => forceLoad,
    () => {
        if (forceLoad) {
            showImage.value = true;
            scrollObserver.disconnect();
        }
    },
    { immediate: true },
);

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
    >
        <img
            v-if="showImage"
            class="content"
            referrerpolicy="same-origin"
            :src="imageURL"
            :width="content.width"
            :height="content.height"
        />

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
</style>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import ScreenCover from "./ScreenCover.vue";
import store from "@/store";

const timeRemaining = ref(3);

onMounted(() => {
    const timer = setInterval(() => {
        timeRemaining.value--;
        if (timeRemaining.value === 0) {
            clearInterval(timer);
        }
    }, 1000);
});

function consent() {
    store.settings.consented = true;
    store.settings.save();

    if (store.settings.searchOnLoad && !store.query.isEmpty()) {
        store.searchPosts();
    }
}

function consentSFW() {
    store.settings.blacklist = store.settings.blacklist.concat([
        {
            name: "rating:explicit",
            type: "unknown",
            count: 0,
        },
        {
            name: "rating:questionable",
            type: "unknown",
            count: 0,
        },
        {
            name: "rating:sensitive",
            type: "unknown",
            count: 0,
        },
    ]);
    consent();
}

function consentNSFW() {
    consent();
}
</script>

<template>
    <ScreenCover blackout />
    <div class="cw-container">
        <div class="content-warning">
            <p>
                Booruview is an image browser for Gelbooru, an image board that
                hosts
                <span class="nsfw-warning">Not Safe For Work</span> content.
            </p>
            <p>
                Do you wish to only view
                <span class="sfw">Safe For Work</span> content?
            </p>
            <p class="btn-container">
                <button
                    class="btn-consent btn-sfw"
                    :disabled="timeRemaining > 0"
                    @click="consentSFW"
                >
                    only view SFW content
                </button>
                <button
                    class="btn-consent btn-nsfw"
                    :disabled="timeRemaining > 0"
                    @click="consentNSFW"
                >
                    view all content (18+)
                </button>
            </p>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-sfw {
    background-color: #70e570;
}

.btn-nsfw {
    background-color: #830b0b;
    color: white;
}

.cw-container {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1;
}

.content-warning {
    padding: 0.8rem;
    max-width: 500px;
}

.nsfw-warning {
    font-weight: bold;
    color: #c00;
}

.sfw {
    font-weight: bold;
    color: #0c0;
}

.btn-consent {
    border-radius: 10px;
    border: none;
    padding: 1rem;
    cursor: pointer;

    &:disabled {
        filter: grayscale(0.5) brightness(0.8);
        cursor: not-allowed;
    }
}

.btn-container {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    justify-content: center;
    max-width: 300px;
    margin: auto;
    margin-top: 3rem;
}

p {
    text-align: center;
}
</style>

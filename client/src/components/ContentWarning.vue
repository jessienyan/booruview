<script setup lang="ts">
import { onMounted, ref } from "vue";
import { type RouteParamsRawGeneric, useRoute, useRouter } from "vue-router";
import { defaultNSFWBlacklist, defaultSFWBlacklist } from "@/blacklist";
import store from "@/store";
import ScreenCover from "./ScreenCover.vue";

const view = ref<"initial" | "nsfw-blacklist">("initial");

const timeRemaining = ref(3);
const router = useRouter();
const route = useRoute();

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
    store.saveSettings();

    // Auto search doesn't trigger on page load if the user hasn't consented.
    // Trigger it once they consent
    let params: RouteParamsRawGeneric;

    if (route.name === "search") {
        params = route.params;
    } else {
        params = { page: 1, query: "" };
    }

    router.push({
        name: "search",
        params,
        force: true,
        replace: true,
    });
}

function consentSFW() {
    store.settings.blacklist = store.settings.blacklist.concat(
        defaultSFWBlacklist(),
    );
    consent();
}

function consentNSFWWithBlacklist() {
    store.settings.blacklist = store.settings.blacklist.concat(
        defaultNSFWBlacklist(),
    );
    consent();
}
</script>

<template>
    <ScreenCover blackout />
    <div class="cw-container">
        <div v-if="view === 'initial'" class="content-warning">
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
                    class="btn-consent btn-green"
                    :disabled="timeRemaining > 0"
                    @click="consentSFW"
                >
                    only view SFW content
                </button>
                <button
                    class="btn-consent btn-red"
                    :disabled="timeRemaining > 0"
                    @click="view = 'nsfw-blacklist'"
                >
                    view all content (18+)
                </button>
            </p>
        </div>

        <div v-else-if="view === 'nsfw-blacklist'" class="content-warning">
            <p>Some content is considered controversial or extreme.</p>
            <p>
                Do you want to use the default blacklist? You can always edit
                your blacklist later.
            </p>
            <p class="btn-container">
                <button
                    class="btn-consent btn-green"
                    :disabled="timeRemaining > 0"
                    @click="consentNSFWWithBlacklist"
                >
                    yes, use default blacklist
                </button>
                <button
                    class="btn-consent btn-red"
                    :disabled="timeRemaining > 0"
                    @click="consent"
                >
                    no, don't filter anything
                </button>
            </p>
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-green {
    background-color: #70e570;
}

.btn-red {
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

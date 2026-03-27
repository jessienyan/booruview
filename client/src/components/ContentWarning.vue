<script setup lang="ts">
import { onMounted, ref } from "vue";
import { type RouteParamsRawGeneric, useRoute, useRouter } from "vue-router";
import { defaultNSFWBlacklist, defaultSFWBlacklist } from "@/blacklist";
import store from "@/store";
import ScreenCover from "./ScreenCover.vue";

const view = ref<"initial" | "nsfw-blacklist">("initial");

const WAIT_TIME_MS = 1800;
const needsToWait = ref(true);
const router = useRouter();
const route = useRoute();

// Add an artificial delay until the buttons are clickable to prevent users from
// accidentally clicking past the consent or accidentally double click.
function makeUserWait() {
    needsToWait.value = true;
    setTimeout(() => {
        needsToWait.value = false;
    }, WAIT_TIME_MS);
}

onMounted(makeUserWait);

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

function showBlacklistOptions() {
    view.value = "nsfw-blacklist";
    makeUserWait();
}

function consentSFW() {
    store.addToBlacklist(defaultSFWBlacklist());
    consent();
}

function consentNSFWWithBlacklist() {
    store.addToBlacklist(defaultNSFWBlacklist());
    consent();
}
</script>

<template>
    <div class="cw-container">
        <ScreenCover blackout />
        <div v-if="view === 'initial'" class="content-warning appear">
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
                    :disabled="needsToWait"
                    @click="consentSFW"
                >
                    only view SFW content
                </button>
                <button
                    class="btn-consent btn-red"
                    :disabled="needsToWait"
                    @click="showBlacklistOptions"
                >
                    view all content (18+)
                </button>
            </p>
        </div>

        <div
            v-else-if="view === 'nsfw-blacklist'"
            class="content-warning appear"
        >
            <p>Some content is considered controversial or extreme.</p>
            <p>
                Do you want to use the default blacklist? You can always edit
                your blacklist later.
            </p>
            <p class="btn-container">
                <button
                    class="btn-consent btn-green"
                    :disabled="needsToWait"
                    @click="consentNSFWWithBlacklist"
                >
                    yes, use default blacklist
                </button>
                <button
                    class="btn-consent btn-red"
                    :disabled="needsToWait"
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

@keyframes fade-out-anim {
    from {
        opacity: 1;
    }
    to {
        opacity: 0;
    }
}

.transition-leave-active {
    animation: 400ms ease fade-out-anim;
}

.btn-green {
    background-color: #70e570;
    color: black;
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
    z-index: 10;
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

    transition: opacity ease 900ms;

    &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
}

@keyframes appear-anim {
    from {
        opacity: 0;
        transform: translateY(12px);
    }
    to {
        opacity: 1;
        transform: translateY(0px);
    }
}
.appear {
    animation: 700ms ease appear-anim;
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

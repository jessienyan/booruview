<script setup lang="tsx">
import { onDeactivated, onUnmounted } from "vue";
import { useRelativeTime } from "@/composable";
import store from "@/store";
import updates from "@/updates";

const relativeTime = useRelativeTime();

function isEntryNew(i: number) {
    const entryNumber = updates.length - i;
    return entryNumber > store.settings.numberUpdatesViewed;
}

const choices = ["wow", "new", "neat", "nice", "cool", "ok"];
const entryStripeText = [];

for (let i = 0; i < updates.length; i++) {
    entryStripeText.push(choices[Math.floor(Math.random() * choices.length)]);
}

function sawUpdates() {
    store.settings.numberUpdatesViewed = updates.length;
    store.saveSettings();
}

// Defer updating this until the user leaves the tab. This way the tab title
// and entries will still appear as new
onUnmounted(sawUpdates);
onDeactivated(sawUpdates);
</script>

<template>
    <div
        v-for="(update, i) in updates"
        class="entry"
        :class="{ 'entry-seen': !isEntryNew(i) }"
    >
        <h2>{{ update.title }}</h2>
        <span v-if="isEntryNew(i)" class="new-stripe">{{
            entryStripeText[i]
        }}</span>
        <component :is="update.component" />
        <p class="timestamp" :title="update.date.toLocaleString()">
            posted {{ relativeTime(update.date) }}
        </p>
    </div>
</template>

<style lang="scss" scoped>
@use "sass:color";
@import "@/assets/colors";

.entry {
    position: relative;
    margin-bottom: 0.8rem;
    padding: 0.8rem;
    overflow: hidden;
    background-color: #2d1b39;

    /* deep pp */
    :deep(p) {
        margin: 0.6rem 0;
    }

    :deep(ul) {
        margin: 0.8rem 0 1rem;
    }

    :deep(li) {
        margin-bottom: 0.4rem;
    }

    p:last-child {
        margin-bottom: 0;
    }

    p.timestamp {
        margin-top: 1rem;
        font-style: italic;
        opacity: 0.7;
        font-size: 0.9em;
    }
}

.entry-seen {
    filter: brightness(0.85) saturate(0.6);
}

.new-stripe {
    position: absolute;
    top: 0px;
    right: 0px;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: #20423f;
    color: #b3b3b3;
    width: 120px;
    height: 24px;
    transform: translate(50%, -50%) rotate(45deg) translateY(170%);
    font-weight: bold;
    font-size: 0.9em;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
}
</style>

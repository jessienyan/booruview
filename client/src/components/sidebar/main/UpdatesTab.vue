<script setup lang="tsx">
import { onMounted, ref } from "vue";
import { useRelativeTime } from "@/composable";
import store from "@/store";
import updates from "@/updates";

const relativeTime = useRelativeTime();
const oldNumberUpdatesViewed = ref(store.settings.numberUpdatesViewed);

onMounted(() => {
    store.settings.numberUpdatesViewed = updates.length;
    store.saveSettings();
});

function isEntryNew(i: number) {
    const entryNumber = updates.length - i;
    return entryNumber > oldNumberUpdatesViewed.value;
}

const choices = ["wow", "neat", "nice", "cool"];
const entryStripeText = [];

for (let i = 0; i < updates.length; i++) {
    entryStripeText.push(choices[Math.floor(Math.random() * choices.length)]);
}
</script>

<template>
    <div
        v-for="(update, i) in updates"
        class="entry"
        :class="[isEntryNew(i) ? 'entry-new' : 'entry-seen']"
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

    background-color: color.scale(
        $color-primary-darker,
        $saturation: 25%,
        $lightness: 2%
    );

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

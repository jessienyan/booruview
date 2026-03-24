<script setup lang="tsx">
import { useRelativeTime } from "@/composable";
import updates from "@/updates";

const relativeTime = useRelativeTime();
</script>

<template>
    <div v-for="update in updates" class="entry entry-new">
        <h2>{{ update.title }}</h2>
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
    margin-bottom: 0.8rem;
    padding: 0.8rem;

    /* the p is deep */
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
    background-color: #221d31;
}

.entry-new {
    background-color: color.scale(
        $color-primary-darker,
        $saturation: 25%,
        $lightness: 2%
    );
}
</style>

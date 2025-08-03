<script setup lang="ts">
import store from "@/store";
import { computed, onMounted, ref, useTemplateRef } from "vue";

const {
    jiggle = false,
    tag,
    fromSearch,
} = defineProps<{
    jiggle?: boolean;
    tag: TagChip;
    fromSearch: boolean;
}>();
const hasJiggled = ref(false);
const chipRef = useTemplateRef("chip");

const cls = computed(() => ({
    [`tag-${tag.tag.type}`]: true,
    strikethrough: tag.style === "strikethrough" || tag.style === "blacklist",
    jiggle: jiggle && !hasJiggled,
}));

onMounted(() => {
    if (jiggle) {
        // Prevents the jiggle animation from playing when the sidebar
        // is opened (display:none triggers animations)
        setTimeout(() => (hasJiggled.value = true), 1000);
    }
});

function onClick() {
    if (!chipRef.value) {
        return;
    }

    // Tag menu is currently open for this chip, close it
    if (store.tagMenu?.ref === chipRef.value) {
        store.tagMenu = null;
        return;
    }

    // HACK: clicking another chip will cause the menu to be closed by ChipMenu.
    // Deferring this allows a new tag to be set after the close is triggered.
    setTimeout(() => {
        store.tagMenu = {
            tag: tag.tag,
            ref: chipRef.value,
            fromSearch,
        };
    }, 0);
}
</script>

<template>
    <div class="chip" :class="cls" ref="chip" @click="onClick">
        <i class="bi bi-check-lg" v-if="tag.style === 'checkmark'"></i>
        <i class="bi bi-ban" v-if="tag.style === 'blacklist'"></i>
        {{ tag.tag.name
        }}<span class="warning" v-if="tag.tag.type === 'deprecated'">
            (deprecated)</span
        >
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

.chip {
    padding: 0.5rem;
    border: none;
    border-radius: 0.5rem;
    display: inline-block;
    word-break: break-all;
    cursor: pointer;
    color: #fff;

    .warning {
        color: #f44;
    }
}

.chip-options {
    button {
        border: 1px solid white;
        background: $color-lightgray;
        color: white;
    }
}

.strikethrough {
    filter: brightness(0.8);
    text-decoration: line-through;
}

@keyframes jiggle-anim {
    0% {
        transform: scale(0.9);
    }

    33% {
        transform: scale(1.1);
    }

    66% {
        transform: scale(0.95);
    }

    100% {
        transform: scale(1);
    }
}

.jiggle {
    animation: 300ms linear 0s jiggle-anim;
}

.tag-deprecated,
.tag-tag {
    background-color: #303030;
    color: hsl(208, 56%, 75%);
}

.tag-artist {
    background-color: #892020;
}

.tag-copyright {
    background-color: #872d87;
}

.tag-character {
    background-color: #1d701d;
}

.tag-metadata {
    background-color: #c98606;
    color: black;
}

.tag-unknown {
    background-color: #6275ae;
    color: #0b1227;
}
</style>

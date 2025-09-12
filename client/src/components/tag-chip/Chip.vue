<script setup lang="ts">
import store from "@/store";
import { computed, onMounted, ref, useTemplateRef } from "vue";
import DropdownMenu from "../DropdownMenu.vue";
import ChipMenuOptions from "./ChipMenuOptions.vue";

const {
    jiggle = false,
    tag,
    canEdit,
} = defineProps<{
    jiggle?: boolean;
    tag: TagChip;
    canEdit: boolean;
}>();
const hasJiggled = ref(false);
const chipRef = useTemplateRef("chip");
const isFavorited = computed(
    () =>
        store.settings.favoriteTags.findIndex(
            (t) => t.name === tag.tag.name,
        ) !== -1,
);
const showOptions = ref(false);

const cls = computed(() => ({
    [`tag-${tag.tag.type}`]: true,
    strikethrough: tag.style === "strikethrough" || tag.style === "blacklist",
    jiggle: jiggle && !hasJiggled,
}));

function onClick() {
    showOptions.value = !showOptions.value;
}

onMounted(() => {
    if (jiggle) {
        // Prevents the jiggle animation from playing when the sidebar
        // is opened (display:none triggers animations)
        setTimeout(() => (hasJiggled.value = true), 1000);
    }
});
</script>

<template>
    <div class="chip" :class="cls" ref="chip" @click="onClick">
        <span class="icons"
            ><i class="bi bi-check-lg" v-if="tag.style === 'checkmark'"></i
            ><i class="bi bi-ban" v-if="tag.style === 'blacklist'"></i
            ><i class="fav-heart bi bi-heart-fill" v-if="isFavorited"></i></span
        >{{ tag.tag.name
        }}<span class="warning" v-if="tag.tag.type === 'deprecated'">
            (deprecated)</span
        >
    </div>

    <DropdownMenu :el="chipRef" v-model:show="showOptions">
        <ChipMenuOptions
            @click="showOptions = false"
            :tag="tag.tag"
            :can-edit="canEdit"
        />
    </DropdownMenu>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";

.chip {
    padding: 0.5rem 0.6rem;
    border: none;
    border-radius: 0.5rem;
    display: inline-block;
    word-break: break-all;
    cursor: pointer;
    color: #fff;
    text-shadow: 0 0 2px black;

    .warning {
        color: #f44;
    }
}

.icons:not(:empty) {
    display: inline-flex;
    justify-content: center;
    align-items: center;
    gap: 4px;
    position: relative;
    top: 1px;
    margin-right: 4px;
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

.fav-heart {
    font-size: 0.9em;
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
    text-shadow: none;
}

.tag-unknown {
    background-color: #6275ae;
    color: #0b1227;
    text-shadow: none;
}
</style>

<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef } from "vue";
import { useRouter } from "vue-router";
import store from "@/store";
import type { ChipActions } from "@/types";
import DropdownMenu from "../DropdownMenu.vue";
import ChipMenuOptions from "./ChipMenuOptions.vue";

const {
	tag,
	actions = {},
	jiggle = false,
	showHeart = true,
} = defineProps<{
	jiggle?: boolean;
	showHeart?: boolean;
	tag: TagChip;
	actions?: ChipActions;
}>();

const hasJiggled = ref(false);
const chipRef = useTemplateRef("chip");
const isFavorited = computed(
	() =>
		store.settings.favoriteTags.findIndex((t) => t.name === tag.tag.name) !==
		-1,
);
const showOptions = ref(false);

const cls = computed(() => ({
	[`tag-${tag.tag.type}`]: true,
	strikethrough: tag.style === "strikethrough" || tag.style === "blacklist",
	jiggle: jiggle && !hasJiggled,
}));

const router = useRouter();
const openInNewTabLink = computed(() => {
	const negated = tag.style === "strikethrough";
	let query = tag.tag.name;
	if (negated) {
		query = "-" + query;
	}

	return router.resolve({
		name: "search",
		params: { page: 1, query },
	}).path;
});

function onClick() {
	if (actions.static) {
		return;
	}

	showOptions.value = !showOptions.value;
}

function onClickLink(e: MouseEvent) {
	if (actions.static) {
		e.preventDefault();
		return;
	}

	// Prevent the link from triggering if the user is just clicking it normally.
	// Control clicking or middle clicking will still trigger the link to open
	if (!e.ctrlKey && e.button === 0) {
		e.preventDefault();
	}
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
    <a :href="openInNewTabLink" target="_blank" @click="onClickLink">
        <div class="chip" :class="cls" ref="chip" @click="onClick">
            <span class="icons"
                ><i class="bi bi-check-lg" v-if="tag.style === 'checkmark'"></i
                ><i class="bi bi-ban" v-if="tag.style === 'blacklist'"></i
                ><i
                    class="fav-heart bi bi-heart-fill"
                    v-if="showHeart && isFavorited"
                ></i></span
            >{{ tag.tag.name
            }}<span class="warning" v-if="tag.tag.type === 'deprecated'">
                (deprecated)</span
            >
        </div>
    </a>

    <DropdownMenu
        v-if="!actions.static"
        :el="chipRef"
        v-model:show="showOptions"
    >
        <ChipMenuOptions
            @click="showOptions = false"
            :tag="tag.tag"
            :actions="actions"
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
    filter: brightness(0.85);
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
    background-color: #7b1414;
}

.tag-copyright {
    background-color: #721d72;
}

.tag-character {
    background-color: #126112;
}

.tag-metadata {
    background-color: #9a6600;
}

.tag-unknown {
    background-color: #44527d;
}
</style>

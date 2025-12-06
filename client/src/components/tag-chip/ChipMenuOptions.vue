<script setup lang="ts">
import { computed, ref } from "vue";
import store from "@/store";
import type { ChipActions } from "@/types";
import { useRouter } from "vue-router";

const showBlacklistConfirm = ref(false);
const {
	onClick,
	tag,
	actions: actionProps = {},
} = defineProps<{
	onClick: () => void;
	tag: Tag;
	actions?: ChipActions;
}>();

const actions = computed(() => ({
	edit: actionProps.edit ?? false,
	blacklist: actionProps.blacklist ?? true,
	includeExcludeRemove: actionProps.includeExcludeRemove ?? true,
	favorite: actionProps.favorite ?? true,
    openInNewTab: actionProps.openInNewTab ?? true,
}));

const isBlacklisted = computed(() => {
	const name = tag.name;
	return store.settings.blacklist.findIndex((t) => t.name === name) !== -1;
});

const isIncluded = computed(() => {
	return store.query.isIncluded(tag.name);
});

const isExcluded = computed(() => {
	return store.query.isExcluded(tag.name);
});

const favoriteIndex = computed(() => {
	return store.settings.favoriteTags.findIndex((t) => tag.name === t.name);
});

const isFavorited = computed(() => favoriteIndex.value !== -1);

const openInNewTabUrl = computed(() => {
    const router = useRouter();
    const url = router.resolve({name: "search", params: {page: 1, query: tag.name}});
    return (new URL(url.path, window.location.origin)).href;
});

function onAdd() {
	store.query.includeTag(tag);
	onClick();
}

function onBlacklist() {
	showBlacklistConfirm.value = true;
}

function onConfirmBlacklist() {
	store.settings.blacklist = store.settings.blacklist.concat(tag);
	store.saveSettings();
	store.query.removeTag(tag);
	onClick();
}

function onExclude() {
	store.query.excludeTag(tag);
	onClick();
}

function onRemove() {
	store.query.removeTag(tag);
	onClick();
}

function onFavorite() {
	store.settings.favoriteTags.push(tag);
	store.saveSettings();
	onClick();
}

function onEdit() {
	store.editTag(tag);
	onClick();
}

function onUnfavorite() {
	store.settings.favoriteTags.splice(favoriteIndex.value, 1);
	store.saveSettings();
	onClick();
}

function onWhitelist() {
	const name = tag.name;
	const i = store.settings.blacklist.findIndex((t) => t.name === name);

	// shouldn't happen
	if (i === -1) {
		return;
	}

	store.settings.blacklist.splice(i, 1);
	store.saveSettings();
	onClick();
}
</script>

<template>
    <button
        class="dropdown-option btn-primary"
        v-if="actions.edit && (isExcluded || isIncluded)"
        @click="onEdit"
    >
        <i class="bi bi-pencil"></i> edit
    </button>
    <button
        class="dropdown-option btn-primary"
        v-if="
            actions.includeExcludeRemove &&
            !isBlacklisted &&
            (isExcluded || !isIncluded)
        "
        @click="onAdd"
    >
        <i class="bi bi-plus-lg"></i> include
    </button>
    <button
        class="dropdown-option btn-primary"
        v-if="
            actions.includeExcludeRemove &&
            !isBlacklisted &&
            (!isExcluded || isIncluded)
        "
        @click="onExclude"
    >
        <i class="bi bi-dash-lg"></i> exclude
    </button>
    <button
        class="dropdown-option btn-primary"
        v-if="
            (actions.includeExcludeRemove && !isBlacklisted && isExcluded) ||
            isIncluded
        "
        @click="onRemove"
    >
        <i class="bi bi-x-lg"></i> remove
    </button>
    <a :href="openInNewTabUrl" target="_blank">
        <button class="dropdown-option btn-primary"
            v-if="actions.openInNewTab && !isBlacklisted && !isExcluded">
            <i class="bi bi-box-arrow-up-right"></i> open in new tab
        </button>
    </a>
    <button
        class="dropdown-option btn-primary"
        v-if="actions.favorite && !isBlacklisted && !isFavorited"
        @click="onFavorite"
    >
        <i class="bi bi-heart"></i> favorite
    </button>
    <button
        class="dropdown-option btn-primary"
        v-if="actions.favorite && !isBlacklisted && isFavorited"
        @click="onUnfavorite"
    >
        <i class="bi bi-heart-fill"></i> unfavorite
    </button>

    <!--
        NOTE: Important to use v-show here to avoid the menu being hidden.
        If the button element is removed from the page the click event won't
        appear to have come from within the container
        -->
    <button
        class="dropdown-option btn-primary blacklist"
        v-show="actions.blacklist && !isBlacklisted && !showBlacklistConfirm"
        @click="onBlacklist"
    >
        <i class="bi bi-ban"></i> blacklist
    </button>
    <button
        class="dropdown-option btn-primary blacklist"
        v-if="!isBlacklisted && showBlacklistConfirm"
        @click="onConfirmBlacklist"
    >
        <i class="bi bi-ban"></i> confirm blacklist
    </button>
    <button
        class="dropdown-option btn-primary rounded"
        v-if="actions.blacklist && isBlacklisted"
        @click="onWhitelist"
    >
        <i class="bi bi-x-lg"></i> remove from blacklist
    </button>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.blacklist {
    color: #ff5d5d;
}
</style>

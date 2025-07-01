<script setup lang="ts">
import store from "@/store";
import { computed, onMounted, ref, useTemplateRef, watch } from "vue";
import { useDismiss } from "@/composable";

const containerRef = useTemplateRef("container");
const showBlacklistConfirm = ref(false);

// Reset confirmation if tag changes
watch([store.tagMenu], () => (showBlacklistConfirm.value = false));
useDismiss([containerRef, store.tagMenu?.ref || null], closeMenu);

const isBlacklisted = computed(() => {
    if (store.tagMenu === null) {
        return false;
    }

    const name = store.tagMenu.tag.name;
    return store.settings.blacklist.findIndex((t) => t.name === name) !== -1;
});

const isIncluded = computed(() => {
    if (store.tagMenu === null) {
        return false;
    }

    return store.query.isIncluded(store.tagMenu.tag.name);
});

const isExcluded = computed(() => {
    if (store.tagMenu === null) {
        return false;
    }

    return store.query.isExcluded(store.tagMenu.tag.name);
});

function closeMenu() {
    store.tagMenu = null;
}

const menuPosition = computed(() => {
    const ref = store.tagMenu?.ref;

    if (ref == null) {
        return;
    }

    const windowHeight = window.innerHeight;
    const allowedMargin = 125;
    const { left, bottom, top } = ref.getBoundingClientRect();

    // Anchor the menu to the bottom of the chip if there is enough space in the viewport
    if (windowHeight - bottom >= allowedMargin) {
        return {
            left: left + "px",
            top: bottom + "px",
        };
    }

    // If space is limited move the anchor to the top of the chip
    const bottomToTop = windowHeight - top;

    return {
        left: left + "px",
        bottom: bottomToTop + "px",
    };
});

function onAdd() {
    if (store.tagMenu === null) {
        return;
    }

    store.query.includeTag(store.tagMenu.tag);
    closeMenu();
}

function onBlacklist() {
    if (store.tagMenu === null) {
        return;
    }

    showBlacklistConfirm.value = true;
}

function onConfirmBlacklist() {
    if (store.tagMenu === null) {
        return;
    }

    store.settings.blacklist = store.settings.blacklist.concat(
        store.tagMenu.tag,
    );
    store.settings.save();
    store.query.removeTag(store.tagMenu.tag);
    closeMenu();
}

function onExclude() {
    if (store.tagMenu === null) {
        return;
    }

    store.query.excludeTag(store.tagMenu.tag);
    closeMenu();
}

function onRemove() {
    if (store.tagMenu === null) {
        return;
    }

    store.query.removeTag(store.tagMenu.tag);
    closeMenu();
}

function onWhitelist() {
    if (store.tagMenu === null) {
        return;
    }

    const name = store.tagMenu.tag.name;
    const i = store.settings.blacklist.findIndex((t) => t.name === name);

    // shouldn't happen
    if (i === -1) {
        return;
    }

    store.settings.blacklist.splice(i, 1);
    store.settings.save();
    closeMenu();
}
</script>

<template>
    <div class="options" ref="container" :style="menuPosition">
        <button
            class="btn-primary option-btn"
            v-if="!isBlacklisted && (isExcluded || !isIncluded)"
            @click="onAdd"
        >
            <i class="bi bi-plus-lg"></i> include
        </button>
        <button
            class="btn-primary option-btn"
            v-if="!isBlacklisted && (!isExcluded || isIncluded)"
            @click="onExclude"
        >
            <i class="bi bi-dash-lg"></i> exclude
        </button>
        <button
            class="btn-primary option-btn"
            v-if="isExcluded || isIncluded"
            @click="onRemove"
        >
            <i class="bi bi-x-lg"></i> remove
        </button>

        <!--
        NOTE: Important to use v-show here to avoid the menu being hidden.
        If the button element is removed from the page the click event won't
        appear to have come from within the container
        -->
        <button
            class="btn-primary option-btn blacklist"
            v-show="!isBlacklisted && !showBlacklistConfirm"
            @click="onBlacklist"
        >
            <i class="bi bi-ban"></i> blacklist
        </button>
        <button
            class="btn-primary option-btn blacklist"
            v-if="!isBlacklisted && showBlacklistConfirm"
            @click="onConfirmBlacklist"
        >
            <i class="bi bi-ban"></i> confirm blacklist
        </button>
        <button
            class="btn-primary option-btn rounded"
            v-if="isBlacklisted"
            @click="onWhitelist"
        >
            <i class="bi bi-x-lg"></i> remove from blacklist
        </button>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.options {
    position: absolute;
    display: flex;
    flex-direction: column;
    z-index: 999;
    width: max-content;
    box-shadow: 0 0 10px black;

    .option-btn {
        padding: 12px;
        text-align: left;

        .bi {
            margin-right: 5px;
        }

        &:first-of-type {
            border-radius: 4px 4px 0 0;
        }

        &:last-of-type {
            border-radius: 0 0 4px 4px;
        }

        &:not(:last-of-type):not(.rounded) {
            border-bottom: none;
        }

        &.rounded {
            border-radius: 4px;
        }
    }
}

.blacklist {
    color: #ff5d5d;
}
</style>

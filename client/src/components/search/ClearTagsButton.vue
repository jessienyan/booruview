<script setup lang="ts">
import { useDismiss, useViewportSize } from "@/composable";
import { computed, ref, useTemplateRef } from "vue";
import store from "@/store";

const btnRef = useTemplateRef("button");
const menuRef = useTemplateRef("menu");
const showMenu = ref(false);
useDismiss([btnRef, menuRef], () => (showMenu.value = false));
const viewport = useViewportSize();

const includeCount = computed(() => store.query._include.size);
const excludeCount = computed(() => store.query._exclude.size);

const menuPosition = computed(() => {
    if (!btnRef.value) {
        return;
    }

    const { width, height } = viewport.value;
    const allowedMargin = 125;
    const { right, bottom, top } = btnRef.value.getBoundingClientRect();
    const rightAlignPos = width - right;

    // Anchor the menu to the bottom of the button if there is enough space in the viewport
    if (height - bottom >= allowedMargin) {
        return {
            right: rightAlignPos + "px",
            top: bottom + "px",
        };
    }

    // If space is limited move the anchor to the top of the button
    const bottomToTop = height - top;

    return {
        right: rightAlignPos + "px",
        bottom: bottomToTop + "px",
    };
});

function clearIncluded() {
    store.query._include.clear();
}

function clearExcluded() {
    store.query._exclude.clear();
}
</script>

<template>
    <button
        ref="button"
        class="btn-gray btn-rounded btn-clear-tags"
        @click="showMenu = !showMenu"
    >
        clear tags
        <i
            class="bi"
            :class="{
                'bi-caret-down-fill': !showMenu,
                'bi-caret-up-fill': showMenu,
            }"
        ></i>
    </button>

    <Teleport to="body">
        <div v-if="showMenu" ref="menu" class="options" :style="menuPosition">
            <button
                v-if="includeCount > 0"
                class="btn-gray option-btn"
                @click="clearIncluded"
            >
                included ({{ includeCount }})
            </button>
            <button
                v-if="excludeCount > 0"
                class="btn-gray option-btn"
                @click="clearExcluded"
            >
                excluded ({{ excludeCount }})
            </button>
            <button
                class="btn-gray option-btn"
                @click="
                    clearIncluded();
                    clearExcluded();
                "
            >
                all
            </button>
        </div>
    </Teleport>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.btn-clear-tags {
    margin-top: 0.8rem;
    margin-left: auto;
    display: block;
}

.options {
    position: absolute;
    display: flex;
    flex-direction: column;
    z-index: 999;
    width: max-content;
    box-shadow: 0 0 0.8rem black;
    margin: 2px 0;

    .option-btn {
        padding: 12px;
        text-align: left;

        .bi {
            margin-right: 0.4rem;
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
</style>

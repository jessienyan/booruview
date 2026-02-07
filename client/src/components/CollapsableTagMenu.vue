<script setup lang="ts">
import { computed, ref, useTemplateRef } from "vue";
import DropdownMenu from "./DropdownMenu.vue";
import Chip from "./tag-chip/Chip.vue";

const { tags } = defineProps<{ tags: TagChip[] }>();

const btnRef = useTemplateRef("btnRef");
const open = ref(false);

const categoryOrder: TagType[] = ["artist", "character", "copyright", "tag", "unknown", "metadata", "deprecated"];

// Tags sorted by category then by name
const sortedTags = computed<TagChip[]>(() =>
	[...tags].sort((a, b) => {
		const category = categoryOrder.indexOf(a.tag.type) - categoryOrder.indexOf(b.tag.type);
		if (category !== 0) {
			return category;
		}
		return a.tag.name.localeCompare(b.tag.name);
	}),
);
</script>

<template>
    <button
        ref="btnRef"
        class="btn-menu-toggle btn-primary btn-rounded"
        @click="open = !open"
    >
        <slot></slot>
        <i
            class="bi"
            :class="{
                'bi-caret-down-fill': !open,
                'bi-caret-up-fill': open,
            }"
        ></i>
    </button>
    <DropdownMenu :el="btnRef" v-model:show="open">
        <div class="tag-list">
            <Chip
                v-for="t of sortedTags"
                :tag="t"
                :actions="{
                    edit: false,
                    favorite: false,
                    includeExcludeRemove: false,
                }"
            />
        </div>
    </DropdownMenu>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.tag-list {
    display: flex;
    flex-direction: column;
    background-color: #1c1c1c;
    gap: 8px;
    padding: 8px;
    max-width: 300px;
    max-height: 350px;
    overflow-y: scroll;

    &:deep(.chip) {
        margin: 0;
        padding: 0.8rem;
    }
}
</style>

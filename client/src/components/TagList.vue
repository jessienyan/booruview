<script setup lang="ts">
import { computed } from "vue";
import Chip from "./tag-chip/Chip.vue";

type listCategories = Record<TagType, TagChip[]>;

const {
    jiggle = false,
    tags = [],
    showCheckmark = false,
} = defineProps<{
    jiggle?: boolean;
    tags: TagChip[];
    showCheckmark?: boolean;
}>();

const categories = computed(() => {
    const ret: listCategories = {
        artist: [],
        character: [],
        copyright: [],
        tag: [],
        metadata: [],
        deprecated: [],
        unknown: [],
    };

    tags.forEach(t => ret[t.tag.type].concat(t));
    // Move deprecated tags into the regular tag section
    ret.tag = ret.tag.concat(ret.deprecated);
    ret.deprecated = [];

    function sortTags(list: TagChip[]) {
        list.sort((a, b) => a.tag.name.localeCompare(b.tag.name));
    }

    sortTags(ret.artist);
    sortTags(ret.character);
    sortTags(ret.copyright);
    sortTags(ret.tag);
    sortTags(ret.metadata);
    sortTags(ret.unknown);

    return ret;
});
</script>

<template>
    <h3 v-if="categories.artist.length > 0">artist</h3>
    <Chip
        v-for="tag in categories.artist"
        :key="tag.tag.name"
        :tag="tag"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.character.length > 0">character</h3>
    <Chip
        v-for="tag in categories.character"
        :key="tag.tag.name"
        :tag="tag"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.copyright.length > 0">copyright</h3>
    <Chip
        v-for="tag in categories.copyright"
        :key="tag.tag.name"
        :tag="tag"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.tag.length > 0">tags</h3>
    <Chip
        v-for="tag in categories.tag"
        :key="tag.tag.name"
        :tag="tag"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.metadata.length > 0">metadata</h3>
    <Chip
        v-for="tag in categories.metadata"
        :key="tag.tag.name"
        :tag="tag"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.unknown.length > 0">raw</h3>
    <Chip
        v-for="tag in categories.unknown"
        :key="tag.tag.name"
        :tag="tag"
        :jiggle="jiggle"
    />
</template>

<style lang="scss" scoped>
h3 {
    margin: 15px 0 10px;
    font-size: 18px;
}
</style>

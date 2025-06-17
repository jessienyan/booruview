<script setup lang="ts">
import { computed } from "vue";
import Chip from "./tag-chip/Chip.vue";

type statefulTag = {
    state: TagState;
    tag: Tag;
};

type listCategories = {
    artist: statefulTag[];
    character: statefulTag[];
    copyright: statefulTag[];
    tag: statefulTag[];
    metadata: statefulTag[];
    unknown: statefulTag[];
};

const {
    jiggle = false,
    includedTags = [],
    excludedTags = [],
    showCheckmark = false,
} = defineProps<{
    jiggle?: boolean;
    includedTags?: Tag[];
    excludedTags?: Tag[];
    showCheckmark?: boolean;
}>();

const categories = computed(() => {
    const ret: listCategories = {
        artist: [],
        character: [],
        copyright: [],
        tag: [],
        metadata: [],
        unknown: [],
    };

    function sortTag(tag: Tag, state: TagState) {
        switch (tag.type) {
            case "artist":
                ret.artist = ret.artist.concat({ state, tag });
                break;
            case "character":
                ret.character = ret.character.concat({ state, tag });
                break;
            case "copyright":
                ret.copyright = ret.copyright.concat({ state, tag });
                break;
            case "deprecated":
            case "tag":
                ret.tag = ret.tag.concat({ state, tag });
                break;
            case "metadata":
                ret.metadata = ret.metadata.concat({ state, tag });
                break;
            case "unknown":
                ret.unknown = ret.unknown.concat({ state, tag });
                break;
        }
    }

    includedTags.forEach((v) => sortTag(v, showCheckmark ? "include" : "none"));
    excludedTags.forEach((v) => sortTag(v, "exclude"));

    ret.artist.sort(({ tag: a }, { tag: b }) => a.name.localeCompare(b.name));
    ret.character.sort(({ tag: a }, { tag: b }) =>
        a.name.localeCompare(b.name),
    );
    ret.copyright.sort(({ tag: a }, { tag: b }) =>
        a.name.localeCompare(b.name),
    );
    ret.tag.sort(({ tag: a }, { tag: b }) => a.name.localeCompare(b.name));
    ret.metadata.sort(({ tag: a }, { tag: b }) => a.name.localeCompare(b.name));
    ret.unknown.sort(({ tag: a }, { tag: b }) => a.name.localeCompare(b.name));

    return ret;
});
</script>

<template>
    <h3 v-if="categories.artist.length > 0">artist</h3>
    <Chip
        v-for="{ state, tag } in categories.artist"
        :key="tag.name"
        :tag="tag"
        :state="state"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.character.length > 0">character</h3>
    <Chip
        v-for="{ state, tag } in categories.character"
        :key="tag.name"
        :tag="tag"
        :state="state"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.copyright.length > 0">copyright</h3>
    <Chip
        v-for="{ state, tag } in categories.copyright"
        :key="tag.name"
        :tag="tag"
        :state="state"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.tag.length > 0">tags</h3>
    <Chip
        v-for="{ state, tag } in categories.tag"
        :key="tag.name"
        :tag="tag"
        :state="state"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.metadata.length > 0">metadata</h3>
    <Chip
        v-for="{ state, tag } in categories.metadata"
        :key="tag.name"
        :tag="tag"
        :state="state"
        :jiggle="jiggle"
    />

    <h3 v-if="categories.unknown.length > 0">raw</h3>
    <Chip
        v-for="{ state, tag } in categories.unknown"
        :key="tag.name"
        :tag="tag"
        :state="state"
        :jiggle="jiggle"
    />
</template>

<style lang="scss" scoped>
h3 {
    margin: 15px 0 10px;
    font-size: 18px;
}
</style>

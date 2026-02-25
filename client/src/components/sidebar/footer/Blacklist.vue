<script setup lang="ts">
import { computed } from "vue";
import { defaultNSFWBlacklist } from "@/blacklist";
import Collapsable from "@/components/Collapsable.vue";
import Chip from "@/components/tag-chip/Chip.vue";
import { useDontShowAgain } from "@/composable";
import store from "@/store";
import { sortTags } from "@/tag";

const blacklist = store.blacklist();
const styledTags = computed<TagChip[]>(() => {
    const sorted = sortTags(blacklist.value);
    const styled = sorted.map<TagChip>((tag) => ({ tag, style: "default" }));
    return styled;
});

const defaultBlacklistVisibility = useDontShowAgain("hide-default-blacklist");

// Consider NSFW enabled if the user hasn't blacklisted rating:explicit
const nsfwEnabled = computed(
    () => blacklist.value.findIndex((t) => t.name === "rating:explicit") === -1,
);

const defaultBlacklistTags = computed(() =>
    defaultNSFWBlacklist().map<TagChip>((tag) => {
        const isBlacklisted =
            blacklist.value.findIndex((t) => t.name === tag.name) !== -1;

        return {
            tag,
            style: isBlacklisted ? "strikethrough" : "default",
        };
    }),
);

function addAllFromDefaultBlacklist() {
    const list: Tag[] = [];

    for (const tag of defaultNSFWBlacklist()) {
        if (blacklist.value.findIndex((t) => t.name === tag.name) === -1) {
            list.push(tag);
        }
    }

    if (list.length) {
        store.setBlacklist(blacklist.value.concat(list));
        store.toast = {
            msg: `added ${list.length} tags to blacklist`,
            type: "info",
        };
    }
}
</script>

<template>
    <div
        v-if="defaultBlacklistVisibility.show.value && nsfwEnabled"
        class="default-blacklist"
    >
        <p>
            A default blacklist is available to filter content which may be
            considered controversial or extreme.
        </p>
        <p>
            <Collapsable text="tags">
                <div class="chips-container">
                    <Chip
                        v-for="t in defaultBlacklistTags"
                        :tag="t"
                        :actions="{
                            edit: false,
                            includeExcludeRemove: false,
                            favorite: false,
                        }"
                    />
                </div>
            </Collapsable>
        </p>
        <p>You can blacklist the tags individually if you prefer.</p>
        <p>
            <button
                class="btn-primary btn-rounded"
                @click="addAllFromDefaultBlacklist"
            >
                add all to blacklist
            </button>
            <button
                class="btn-primary-darker btn-rounded"
                @click="defaultBlacklistVisibility.onHide"
            >
                don't show again
            </button>
        </p>
    </div>

    <p v-if="styledTags.length > 0">
        <Collapsable text="view blacklist">
            <div class="chips-container">
                <Chip v-for="t in styledTags" :tag="t" />
            </div>
        </Collapsable>
    </p>

    <template v-else>
        <p>
            Any tags you've blacklisted will appear here. Posts that have
            blacklisted tags will be hidden.
        </p>
        <p>To blacklist, click a tag and then select "blacklist".</p>
    </template>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";

.default-blacklist {
    padding: 0.8rem;
    background-color: $color-primary-darker;

    p {
        color: $color-primary-light;

        &:first-child {
            margin-top: 0;
        }

        &:last-child {
            margin-bottom: 0;
        }
    }
}

.chips-container {
    margin-top: 1rem;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 0.3rem;
}
</style>

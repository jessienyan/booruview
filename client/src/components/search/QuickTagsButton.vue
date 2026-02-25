<script setup lang="ts">
import { computed, ref, useTemplateRef } from "vue";
import store from "@/store";
import { sortTags } from "@/tag";
import DropdownMenu from "../DropdownMenu.vue";
import Chip from "../tag-chip/Chip.vue";

const blacklist = store.blacklist();
const styledTag = (tag: Tag) =>
    computed<TagChip>(() => {
        let style: ChipStyle = "default";

        if (store.query.isIncluded(tag.name)) {
            style = "checkmark";
        } else if (store.query.isExcluded(tag.name)) {
            style = "strikethrough";
        } else {
            const isBlacklisted =
                blacklist.value.findIndex((t) => t.name === tag.name) !== -1;
            if (isBlacklisted) {
                style = "strikethrough";
            }
        }

        return { tag, style };
    });

interface Examples {
    score: TagChip[];
    sorting: TagChip[];
    rating: TagChip[];
}
const examples = computed<Examples>(() => ({
    score: [
        styledTag({ count: 0, name: "score:>=100", type: "unknown" }).value,
        styledTag({ count: 0, name: "score:>5", type: "unknown" }).value,
        styledTag({ count: 0, name: "score:<20", type: "unknown" }).value,
    ],
    sorting: [
        styledTag({ count: 0, name: "sort:score", type: "unknown" }).value,
        styledTag({ count: 0, name: "sort:random", type: "unknown" }).value,
    ],
    rating: [
        styledTag({ count: 0, name: "rating:general", type: "metadata" }).value,
        styledTag({ count: 0, name: "rating:sensitive", type: "metadata" })
            .value,
        styledTag({ count: 0, name: "rating:questionable", type: "metadata" })
            .value,
        styledTag({ count: 0, name: "rating:explicit", type: "metadata" })
            .value,
    ],
}));

type Tab = "favs" | "sortfilter";
const currentTab = ref<Tab>("favs");

const open = ref(false);
const btnRef = useTemplateRef("button");

const favTags = store.favoriteTags();

// Favorited tags sorted by category then by name
const sortedFavTags = computed<TagChip[]>(() => {
    const sorted = sortTags(favTags.value);
    const styled: TagChip[] = [];

    for (const tag of sorted) {
        let style: ChipStyle = "default";

        if (store.query.isIncluded(tag.name)) {
            style = "checkmark";
        } else if (store.query.isExcluded(tag.name)) {
            style = "strikethrough";
        }

        styled.push({
            tag,
            style,
        });
    }

    return styled;
});
</script>

<template>
    <button
        ref="button"
        class="btn-quick-tags btn-primary"
        @click="open = !open"
    >
        <i class="bi bi-tags-fill"></i>{{ " "
        }}<i
            class="bi"
            :class="{
                'bi-caret-down-fill': !open,
                'bi-caret-up-fill': open,
            }"
        ></i>
    </button>
    <DropdownMenu :el="btnRef" v-model:show="open">
        <div class="tabs">
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'favs' }"
                @click="currentTab = 'favs'"
            >
                favs
            </button>
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'sortfilter' }"
                @click="currentTab = 'sortfilter'"
            >
                sort / filter
            </button>
        </div>

        <div class="content-container">
            <template v-if="currentTab === 'favs'">
                <div
                    v-if="sortedFavTags.length > 0"
                    class="chip-list chip-list-vertical"
                >
                    <Chip
                        v-for="t of sortedFavTags"
                        :tag="t"
                        :show-heart="false"
                    />
                </div>
                <p v-else>You don't have any favorited tags.</p>
            </template>

            <template v-else>
                <p>
                    These are some examples, check the help tab for more info.
                </p>

                <h3>sorting</h3>

                <div class="chip-list">
                    <Chip v-for="t of examples.sorting" :tag="t" />
                </div>

                <h3>score</h3>

                <div class="chip-list">
                    <Chip
                        v-for="t of examples.score"
                        :tag="t"
                        :actions="{ edit: true }"
                    />
                </div>

                <p>Tip: include and then edit the tag by clicking it again.</p>

                <h3>content rating</h3>

                <div class="chip-list">
                    <Chip v-for="t of examples.rating" :tag="t" />
                </div>
            </template>
        </div>
    </DropdownMenu>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";

.tabs {
    display: flex;
    gap: 0.4rem;
    padding: 0.8rem;
    padding-bottom: 0;
    background-color: #151515;
    border-bottom: 1px solid #555;
}

.tab-btn {
    border: 1px solid #555;
    border-bottom: none;
    padding: 0.4rem 0.8rem;
    border-radius: 4px 4px 0 0;
    background-color: #1e1e1e;
    color: #999;
    cursor: pointer;
    width: min-content;
    text-wrap: nowrap;

    &.active {
        border-color: #695675;
        color: #bb9fce;
        background-color: #342b3a;
    }
}

.btn-quick-tags {
    border-top-right-radius: 4px;
    border-bottom-right-radius: 4px;
}

.content-container {
    padding: 0.8rem;
    background-color: #1c1c1c;
    overflow-y: scroll;
    max-width: 300px;
    max-height: 350px;

    p:first-child {
        margin-top: 0;
    }

    p:last-child {
        margin-bottom: 0;
    }
}

.chip-list {
    display: flex;
    flex-direction: column;
    gap: 0.8rem;

    &:deep(.chip) {
        margin: 0;
        text-wrap: nowrap;
    }
}
</style>

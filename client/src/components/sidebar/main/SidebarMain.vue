<script setup lang="ts">
import { type Component, computed, ref } from "vue";
import news from "@/news";
import store from "@/store";
import NewsTab from "./NewsTab.vue";
import RecentTab from "./RecentTab.vue";
import SavedTab from "./SavedTab.vue";
import SearchTab from "./SearchTab.vue";

type Tab = "search" | "recent" | "saved" | "news";
const tabComponents: Record<Tab, Component> = {
    search: SearchTab,
    recent: RecentTab,
    saved: SavedTab,
    news: NewsTab,
};
const currentTab = ref<Tab>("search");

function switchTab(tab: Tab) {
    currentTab.value = tab;
}

const numUnreadNews = computed(() => {
    let i = 0;
    for (const { date } of news) {
        if (date > store.settings.newsLastViewedAt) {
            i++;
        }
    }
    return i;
});
</script>

<template>
    <div class="tab-container">
        <header class="tabs">
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'search' }"
                @click="switchTab('search')"
            >
                search
            </button>
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'recent' }"
                @click="switchTab('recent')"
            >
                recent
            </button>
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'saved' }"
                @click="switchTab('saved')"
            >
                saved
            </button>

            <div class="spacer" />

            <button
                class="tab-btn"
                :class="{
                    active: currentTab === 'news',
                    highlight: numUnreadNews > 0,
                }"
                @click="switchTab('news')"
            >
                news
                <template v-if="numUnreadNews > 0"
                    >({{ numUnreadNews }})</template
                >
            </button>
        </header>

        <div class="tab-content-container">
            <div class="tab-content">
                <KeepAlive>
                    <component :is="tabComponents[currentTab]" />
                </KeepAlive>
            </div>
        </div>
    </div>
</template>

<style scoped lang="scss">
.tab-container {
    flex: 1;
    display: flex;
    flex-direction: column;

    // Enough height to show the search box and search button
    min-height: 150px;

    &:not(.closed) {
        .tabs {
            border-bottom: 1px solid #555;
        }
    }
}

.tabs {
    display: flex;
    padding: 0 0.8rem;
    gap: 0.4rem;

    .spacer {
        flex: 1;
    }
}

.tab-btn {
    border: 1px solid #555;
    border-bottom: none;
    padding: 0.4rem 0.8rem;
    border-radius: 4px 4px 0 0;
    background-color: #1e1e1e;
    color: #999;
    cursor: pointer;
    font-size: 16px;
    max-width: 100%;

    &.active {
        border-color: #695675;
        color: #bb9fce;
        background-color: #342b3a;
    }

    &.highlight {
        font-weight: bold;
    }
}

.close-btn {
    margin-left: auto;
}

.tab-content-container {
    flex: 1;
    min-height: 0;
    overflow-y: scroll;
}

.tab-content {
    margin: 0.8rem;
}

p,
li {
    font-size: 16px;
    color: #999;

    a,
    a:visited {
        color: #bb9fce;
    }
}

select {
    background-color: #222;
    border: 1px solid #555;
    color: #999;
    padding: 0.4rem;
    border-radius: 0.4rem;
}
</style>

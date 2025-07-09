<script setup lang="ts">
import { ref } from "vue";
import SearchTab from "./SearchTab.vue";
import RecentTab from "./RecentTab.vue";
import NewFeature from "@/components/NewFeature.vue";
import { useNewFeatureIndicator } from "@/composable";

type Tab = "search" | "recent";
const currentTab = ref<Tab>("search");
const featRecent = useNewFeatureIndicator("recent_tab", new Date("2025-07-08"));

function switchTab(tab: Tab) {
    currentTab.value = tab;
}
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
                @click="
                    switchTab('recent');
                    featRecent.onSeen();
                "
            >
                recent
                <NewFeature v-if="featRecent.show.value" />
            </button>
        </header>

        <div class="tab-content-container">
            <div class="tab-content">
                <KeepAlive>
                    <SearchTab v-if="currentTab === 'search'" />
                    <RecentTab v-else-if="currentTab === 'recent'" />
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

    // 140px is enough height to show the search box and search button
    min-height: 140px;

    &:not(.closed) {
        .tabs {
            border-bottom: 1px solid #555;
        }
    }
}

.tabs {
    display: flex;
    padding: 0 10px;
    gap: 5px;
}

.tab-btn {
    border: 1px solid #555;
    border-bottom: none;
    padding: 5px 10px;
    border-radius: 4px 4px 0 0;
    background-color: #1e1e1e;
    color: #999;
    cursor: pointer;
    font-size: 16px;

    &.active {
        border-color: #695675;
        color: #bb9fce;
        background-color: #342b3a;
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
    margin: 1rem;
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
    padding: 5px;
    border-radius: 5px;
}
</style>

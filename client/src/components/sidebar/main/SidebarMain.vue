<script setup lang="ts">
import { ref } from "vue";
import SearchTab from "./SearchTab.vue";
import RecentTab from "./RecentTab.vue";

type Tab = "search" | "recent";
const currentTab = ref<Tab>("search");

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
                @click="switchTab('recent')"
            >
                recent
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
    padding: 0 0.8rem;
    gap: 0.4rem;
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
    width: min-content;

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

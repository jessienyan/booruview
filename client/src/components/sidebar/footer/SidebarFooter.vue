<script setup lang="ts">
import { computed, ref } from "vue";
import AboutTab from "./AboutTab.vue";
import HelpTab from "./HelpTab.vue";
import SettingsTab from "./SettingsTab.vue";
import store from "@/store";
import BlacklistTab from "./BlacklistTab.vue";

type Tab = "about" | "help" | "settings" | "blacklist";
const currentTab = ref<Tab>("about");

// CBA writing the full setting each time
const closed = computed({
    get: () => store.settings.sidebarTabsHidden,
    set: (val: boolean) => (store.settings.sidebarTabsHidden = val),
});

function switchTab(tab: Tab) {
    currentTab.value = tab;
    closed.value = false;
    store.saveSettings();
}

function toggleClose() {
    closed.value = !closed.value;
    store.saveSettings();
}
</script>

<template>
    <div class="tab-container" :class="{ closed }">
        <header class="tabs">
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'about' && !closed }"
                @click="switchTab('about')"
            >
                about
            </button>
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'help' && !closed }"
                @click="switchTab('help')"
            >
                help
            </button>
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'settings' && !closed }"
                @click="switchTab('settings')"
            >
                settings
            </button>
            <button
                class="tab-btn"
                :class="{ active: currentTab === 'blacklist' && !closed }"
                @click="switchTab('blacklist')"
            >
                blacklist
            </button>

            <button
                v-if="!closed"
                class="tab-btn close-btn"
                @click="toggleClose"
            >
                <i class="bi bi-chevron-down"></i>
            </button>
        </header>

        <div class="tab-content-container" v-if="!closed">
            <div class="tab-content">
                <KeepAlive>
                    <HelpTab v-if="currentTab === 'help'" />
                    <AboutTab v-else-if="currentTab === 'about'" />
                    <SettingsTab v-else-if="currentTab === 'settings'" />
                    <BlacklistTab v-else-if="currentTab === 'blacklist'" />
                </KeepAlive>
            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.tab-container {
    display: flex;
    flex-direction: column;
    min-height: 0;

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
</style>

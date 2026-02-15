<script setup lang="ts">
import { type Component, ref } from "vue";
import store from "@/store";
import AboutTab from "./AboutTab.vue";
import AccountTab from "./AccountTab.vue";
import HelpTab from "./HelpTab.vue";
import SettingsTab from "./SettingsTab.vue";

type Tab = "about" | "help" | "settings" | "account";
const tabComponents: Record<Tab, Component> = {
    about: AboutTab,
    help: HelpTab,
    settings: SettingsTab,
    account: AccountTab,
};
const currentTab = ref<Tab>("about");

function switchTab(tab: Tab) {
    currentTab.value = tab;
    store.settings.sidebarTabsHidden = false;
    store.saveSettings();
}

function toggleClose() {
    store.settings.sidebarTabsHidden = !store.settings.sidebarTabsHidden;
    store.saveSettings();
}
</script>

<template>
    <div
        class="tab-container"
        :class="{ closed: store.settings.sidebarTabsHidden }"
    >
        <header class="tabs">
            <button
                class="tab-btn"
                :class="{
                    active:
                        currentTab === 'about' &&
                        !store.settings.sidebarTabsHidden,
                }"
                @click="switchTab('about')"
            >
                about
            </button>
            <button
                class="tab-btn"
                :class="{
                    active:
                        currentTab === 'help' &&
                        !store.settings.sidebarTabsHidden,
                }"
                @click="switchTab('help')"
            >
                help
            </button>
            <button
                class="tab-btn"
                :class="{
                    active:
                        currentTab === 'settings' &&
                        !store.settings.sidebarTabsHidden,
                }"
                @click="switchTab('settings')"
            >
                settings
            </button>

            <button
                class="tab-btn"
                :class="{
                    active:
                        currentTab === 'account' &&
                        !store.settings.sidebarTabsHidden,
                }"
                @click="switchTab('account')"
            >
                account
            </button>

            <button
                v-if="!store.settings.sidebarTabsHidden"
                class="tab-btn close-btn"
                @click="toggleClose"
            >
                <i class="bi bi-chevron-down"></i>
            </button>
        </header>

        <div
            class="tab-content-container"
            v-if="!store.settings.sidebarTabsHidden"
        >
            <div class="tab-content">
                <KeepAlive>
                    <component :is="tabComponents[currentTab]" />
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

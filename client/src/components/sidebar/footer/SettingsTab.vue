<script setup lang="ts">
import store, { type ColumnSizing } from "@/store";

const columnSizingOptions: Record<ColumnSizing, string> = {
    dynamic: "dynamic",
    fixed: "fixed",
};

function onChangeColCount(e: Event) {
    store.settings.columnCount = parseInt((e.target as HTMLInputElement).value);
    store.settings.save();
}

function onChangeColSizing(e: Event) {
    store.settings.columnSizing = (e.target as HTMLInputElement).value as any;
    store.settings.save();
}

function onChangeColWidth(e: Event) {
    store.settings.columnWidth = parseInt((e.target as HTMLInputElement).value);
    store.settings.save();
}

function onChangeCloseSidebarOnSearch(e: Event) {
    store.settings.closeSidebarOnSearch = (
        e.target as HTMLInputElement
    ).checked;
    store.settings.save();
}

function onChangeSearchOnLoad(e: Event) {
    store.settings.searchOnLoad = (e.target as HTMLInputElement).checked;
    store.settings.save();
}

function onChangeHighResImages(e: Event) {
    store.settings.highResImages = (e.target as HTMLInputElement).checked;
    store.settings.save();
}

function onChangeAutoplayVideos(e: Event) {
    store.settings.autoplayVideo = (e.target as HTMLInputElement).checked;
    store.settings.save();
}

function onChangeMuteVideos(e: Event) {
    store.settings.muteVideo = (e.target as HTMLInputElement).checked;
    store.settings.save();
}
</script>

<template>
    <div class="settings-container">
        <h3>layout</h3>

        <div class="input-group">
            <label># of columns</label>
            <div class="input">
                <select
                    :value="store.settings.columnSizing"
                    @change="onChangeColSizing"
                >
                    <option
                        v-for="(label, val) in columnSizingOptions"
                        :value="val"
                    >
                        {{ label }}
                    </option>
                </select>
            </div>
        </div>

        <div class="input-group" v-if="store.settings.columnSizing === 'fixed'">
            <label>column count</label>
            <div class="input">
                <input
                    type="range"
                    min="1"
                    max="20"
                    step="1"
                    :value="store.settings.columnCount"
                    @input="onChangeColCount"
                />
                <span class="value">{{ store.settings.columnCount }}</span>
            </div>
        </div>

        <div
            class="input-group"
            v-if="store.settings.columnSizing === 'dynamic'"
        >
            <label>max column width</label>
            <div class="input">
                <input
                    type="range"
                    min="100"
                    max="1000"
                    step="10"
                    :value="store.settings.columnWidth"
                    @input="onChangeColWidth"
                />
                <span class="value">{{ store.settings.columnWidth }}px</span>
            </div>
        </div>

        <h3>search</h3>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.searchOnLoad"
                    @change="onChangeSearchOnLoad"
                />
                auto-search when page loads if search isn't empty</label
            >
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.closeSidebarOnSearch"
                    @change="onChangeCloseSidebarOnSearch"
                />
                searching closes sidebar</label
            >
        </div>

        <h3>content</h3>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.highResImages"
                    @change="onChangeHighResImages"
                />
                high resolution images (uncheck if slow connection)</label
            >
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.autoplayVideo"
                    @change="onChangeAutoplayVideos"
                />
                autoplay videos</label
            >
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.muteVideo"
                    @change="onChangeMuteVideos"
                />
                mute videos</label
            >
        </div>
    </div>
</template>

<style lang="scss" scoped>
.input {
    display: flex;

    input,
    select {
        flex: 1;
    }

    .value {
        min-width: 60px;
        text-align: center;
    }
}

.input-group {
    margin: 1rem 0 1rem 1rem;
}

label + .input {
    margin-top: 0.4rem;
}

label {
    display: block;
}
</style>

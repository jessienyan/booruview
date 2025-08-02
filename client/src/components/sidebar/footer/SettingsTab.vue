<script setup lang="ts">
import store, {
    type ColumnSizing,
    type FullscreenViewMenuAnchorPoint,
    type PageLoadAutoSearch,
} from "@/store";

const columnSizingOptions: Record<ColumnSizing, string> = {
    dynamic: "dynamic",
    fixed: "fixed",
};

const fullscreenViewMenuAnchorOptions: Record<
    FullscreenViewMenuAnchorPoint,
    string
> = {
    topleft: "top left",
    topcenter: "top center",
    topright: "top right",
    right: "right",
    bottomright: "bottom right",
    bottomcenter: "bottom center",
    bottomleft: "bottom left",
    left: "left",
};

const searchOnPageLoadOptions: Record<PageLoadAutoSearch, string> = {
    always: "always",
    "if-query": "if search isn't empty",
    never: "never",
};

function onChangeColCount(e: Event) {
    store.settings.columnCount = parseInt((e.target as HTMLInputElement).value);
    store.saveSettings();
}

function onChangeColSizing(e: Event) {
    store.settings.columnSizing = (e.target as HTMLInputElement)
        .value as ColumnSizing;
    store.saveSettings();
}

function onChangeColWidth(e: Event) {
    store.settings.columnWidth = parseInt((e.target as HTMLInputElement).value);
    store.saveSettings();
}

function onChangeFullscreenViewMenuAnchor(e: Event) {
    store.settings.fullscreenViewMenuAnchor = (e.target as HTMLInputElement)
        .value as FullscreenViewMenuAnchorPoint;
    store.saveSettings();
}

function onChangeFullscreenViewMenuRotate(e: Event) {
    store.settings.fullscreenViewMenuRotate = (
        e.target as HTMLInputElement
    ).checked;
    store.saveSettings();
}

function onChangeCloseSidebarOnSearch(e: Event) {
    store.settings.closeSidebarOnSearch = (
        e.target as HTMLInputElement
    ).checked;
    store.saveSettings();
}

function onChangeSearchOnPageLoad(e: Event) {
    store.settings.searchOnPageLoad = (e.target as HTMLInputElement)
        .value as any;
    store.saveSettings();
}

function onChangeHighResImages(e: Event) {
    store.settings.highResImages = (e.target as HTMLInputElement).checked;
    store.saveSettings();
}

function onChangeAutoplayVideos(e: Event) {
    store.settings.autoplayVideo = (e.target as HTMLInputElement).checked;
    store.saveSettings();
}

function onChangeMuteVideos(e: Event) {
    store.settings.muteVideo = (e.target as HTMLInputElement).checked;
    store.saveSettings();
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

        <div class="input-group">
            <label>fullscreen view menu position</label>
            <div class="input">
                <select
                    :value="store.settings.fullscreenViewMenuAnchor"
                    @change="onChangeFullscreenViewMenuAnchor"
                >
                    <option
                        v-for="(label, val) in fullscreenViewMenuAnchorOptions"
                        :value="val"
                    >
                        {{ label }}
                    </option>
                </select>
            </div>
        </div>

        <div class="input-group">
            <label>
                <input
                    type="checkbox"
                    :checked="store.settings.fullscreenViewMenuRotate"
                    @change="onChangeFullscreenViewMenuRotate"
                />
                rotate menu 90 degrees</label
            >
        </div>

        <h3>search</h3>

        <div class="input-group">
            <label>auto-search when page loads</label>
            <div class="input">
                <select
                    :value="store.settings.searchOnPageLoad"
                    @change="onChangeSearchOnPageLoad"
                >
                    <option
                        v-for="(label, val) in searchOnPageLoadOptions"
                        :value="val"
                    >
                        {{ label }}
                    </option>
                </select>
            </div>
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

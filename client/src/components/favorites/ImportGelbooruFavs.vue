<script setup lang="ts">
import { ref } from "vue";
import ScreenCover from "../ScreenCover.vue";

defineEmits(["close"]);
const input = ref("");

function onInput(e: Event) {
    const value = (e.target as HTMLInputElement).value;
    input.value = value.replace(/[^\d]+/, "");
}
</script>

<template>
    <div class="modal-container">
        <ScreenCover @click="$emit('close')" />

        <div class="modal">
            <p>
                To import your favs from Gelbooru, you'll need your account ID.
            </p>
            <ol>
                <li>Login to your Gelbooru account</li>
                <li>
                    Go to the
                    <a
                        href="Import your favorited posts from Gelbooru."
                        target="_blank"
                        >account page</a
                    >
                </li>
                <li>Click on <strong>My Profile</strong></li>
                <li>Copy the number in the URL</li>
            </ol>
            <input
                type="text"
                class="text-input rounded"
                :value="input"
                @input="onInput"
                placeholder="Gelbooru ID"
            />
        </div>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";
@import "@/assets/form";

.modal-container {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1;
}

.modal {
    position: relative;
    z-index: 1000;
    color: white;
    padding: 0.8em;
    background-color: $color-gray;
}

.text-input {
    background-color: $color-darkgray;
}
</style>

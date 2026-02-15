<script setup lang="ts">
import { ref } from "vue";

const username = ref<string | null>(null);
const password = ref<string | null>(null);

function onSubmit() {}
</script>

<template>
    <form @submit="onSubmit">
        <input
            v-model="username"
            class="text-input rounded"
            type="text"
            placeholder="username"
            maxlength="16"
            minLength="3"
            required
        />
        <input
            v-model="password"
            class="text-input rounded"
            type="password"
            placeholder="password"
            required
        />

        <button class="submit btn-primary btn-rounded" type="submit">
            login
        </button>
    </form>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/form";

form {
    padding: 0 1em;
}

.text-input {
    width: 100%;
    margin-bottom: 1em;
}

.submit {
    width: 100%;
    display: block;
}
</style>

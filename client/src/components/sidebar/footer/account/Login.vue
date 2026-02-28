<script setup lang="ts">
import { ref } from "vue";
import store from "@/store";

const username = ref("");
const password = ref("");
const submitting = ref(false);

function onSubmit() {
    if (submitting.value) {
        return;
    }

    submitting.value = true;
    store.login(username.value, password.value).finally(() => {
        submitting.value = false;
    });
}
</script>

<template>
    <form @submit.prevent="onSubmit">
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

        <button
            class="submit btn-primary btn-rounded"
            type="submit"
            :disabled="submitting"
            @click="onSubmit"
        >
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

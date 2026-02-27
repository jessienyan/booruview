<script lang="ts" setup>
import { computed, ref, watch } from "vue";
import store from "@/store";

const showForm = ref(false);
const currentPassword = ref("");
const newPassword = ref("");
const newPassword2 = ref("");

// Reset form when it's closed
watch([showForm], () => {
    if (!showForm.value) {
        currentPassword.value = "";
        newPassword.value = "";
        newPassword2.value = "";
    }
});

const canSubmit = computed(
    () =>
        store.account &&
        currentPassword.value &&
        newPassword.value &&
        newPassword.value === newPassword2.value,
);

async function submit() {
    if (!canSubmit.value) {
        return;
    }

    try {
        const body = {
            current_password: currentPassword.value,
            new_password: newPassword.value,
        };

        const resp = await fetch("/api/account/password", {
            method: "POST",
            body: JSON.stringify(body),
            headers: {
                Authorization: `Bearer ${store.account!.authToken}`,
                "Content-Type": "application/json",
            },
        });

        const data = await resp.json();
        if (data.error) {
            store.toast = {
                msg: data.error,
                type: "error",
            };
            return;
        }

        store.account!.authToken = data.auth_token;
        store.saveAccountCredentials();
        store.toast = {
            msg: "Password changed",
            type: "info",
        };
        showForm.value = false;
    } catch (e) {
        console.log(e);
        store.toast = {
            msg: "Something went wrong :(",
            type: "error",
        };
    }
}
</script>

<template>
    <template v-if="store.account">
        <button
            class="btn-primary btn-rounded btn-block"
            @click="showForm = true"
            :disabled="showForm"
        >
            Change Password
        </button>
        <div v-if="showForm" class="form-container">
            <p>
                You will need to login again on your other devices after
                changing your password.
            </p>
            <form @submit.prevent="submit">
                <input
                    v-model="currentPassword"
                    class="text-input input-block rounded"
                    type="password"
                    placeholder="current password"
                    required
                />
                <input
                    v-model="newPassword"
                    class="text-input input-block rounded"
                    type="password"
                    placeholder="new password"
                    required
                />
                <input
                    v-model="newPassword2"
                    class="text-input input-block rounded"
                    type="password"
                    placeholder="confirm password"
                    required
                />
            </form>
            <p>
                <button
                    class="btn-primary btn-rounded btn-block"
                    :disabled="!canSubmit"
                    @click="submit"
                >
                    Submit</button
                ><button
                    class="btn-gray btn-rounded btn-block"
                    @click="showForm = false"
                >
                    Cancel
                </button>
            </p>
        </div>
    </template>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

.form-container {
    background-color: $color-darkgray;
    padding: 0.1px 1em;

    p {
        margin: 1.5em 0;

        &:last-child {
            margin-bottom: inherit;
        }
    }
}
</style>

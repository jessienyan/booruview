<script setup lang="ts">
import { ref } from "vue";
import Collapsable from "@/components/Collapsable.vue";
import store from "@/store";
import DeleteAccount from "./account/DeleteAccount.vue";
import FAQ from "./account/FAQ.vue";
import Login from "./account/Login.vue";
import Register from "./account/Register.vue";

const showRegisterForm = ref(false);

function logout() {
    store.account = null;
    store.saveAccountCredentials();
}
</script>

<template>
    <template v-if="!store.account">
        <template v-if="showRegisterForm">
            <h3>register</h3>
            <Register />
            <p>
                <a
                    class="change-form"
                    href="#"
                    @click.prevent="showRegisterForm = false"
                    >I already have an account</a
                >
            </p>
        </template>
        <template v-else>
            <h3>login</h3>
            <Login />
            <p>
                <a
                    class="change-form"
                    href="#"
                    @click.prevent="showRegisterForm = true"
                    >Create an account</a
                >
            </p>
        </template>
    </template>
    <template v-else>
        <p v-if="store.account">
            <strong>You are signed in as {{ store.account.username }}.</strong>
        </p>
        <button class="btn-primary btn-rounded btn-block" @click="logout">
            Logout
        </button>
        <DeleteAccount />
    </template>
    <FAQ />
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

.collapsable {
    background-color: $color-primary-darker;
    color: $color-primary-light;
    padding: 0.1px 1em;
}

.text-input {
    text-align: center;
}

.change-form {
    display: block;
    text-align: center;
}
</style>

<script setup lang="ts">
import { ref } from "vue";
import Collapsable from "@/components/Collapsable.vue";
import store from "@/store";
import DeleteAccount from "./DeleteAccount.vue";
import Login from "./Login.vue";
import Register from "./Register.vue";

const showRegisterForm = ref(false);

function logout() {
    store.account = null;
    store.saveAccountCredentials();
    store.toast = {
        msg: "You have been logged out",
        type: "info",
    };
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
    <div class="faq">
        <Collapsable text="Do I have to register an account?" :button="false">
            <div class="collapsable">
                <p>No, accounts are entirely optional.</p>
            </div>
        </Collapsable>

        <Collapsable text="What data is stored?" :button="false">
            <div class="collapsable">
                <p>
                    Accounts store your favorite tags and posts, blacklist, and
                    search history on Booruview.
                </p>
            </div>
        </Collapsable>
        <Collapsable text="How do accounts work?" :button="false">
            <div class="collapsable">
                <p>
                    It lets you access your data from multiple devices and keeps
                    those devices in sync.
                </p>
                <p>
                    Without an account, your data is only saved in your browser.
                    You're at risk of losing your data if your device breaks or
                    you accidentally clear your browser data.
                </p>
            </div>
        </Collapsable>

        <Collapsable text="How is my data synced?" :button="false">
            <div class="collapsable">
                <p>
                    While logged in, your data is loaded when you open Booruview
                    or refresh the page. Any changes you make are sent back to
                    Booruview.
                </p>
            </div>
        </Collapsable>
    </div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

.faq {
    margin-top: 2em;
}

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

<script setup lang="ts">
import Collapsable from "@/components/Collapsable.vue";
</script>

<template>
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
@import "@/assets/colors";

.faq {
    margin-top: 2em;
}

.collapsable {
    background-color: $color-primary-darker;
    color: $color-primary-light;
    padding: 0.1px 1em;
}
</style>

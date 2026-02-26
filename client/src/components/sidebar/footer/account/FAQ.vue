<script setup lang="ts">
import Collapsable from "@/components/Collapsable.vue";
</script>

<template>
    <div class="faq">
        <Collapsable text="What is an account?" :button="false">
            <div class="collapsable">
                <p>
                    Accounts lets you easily sync your data across your devices.
                    They're free, entirely optional, and no email required.
                </p>
                <p>
                    Accounts store your favorite posts, favorite tags,
                    blacklist, and search history.
                </p>
                <p>
                    Things like your layout settings are not synced, however, in
                    case you want different settings for different devices.
                </p>
            </div>
        </Collapsable>

        <Collapsable
            text="Can I sync up my device with my account?"
            :button="false"
        >
            <div class="collapsable">
                <p>
                    Yes. Login to the device that has your existing favorites,
                    then click the <strong>Upload Data</strong> button in the
                    account tab. This will merge your data from your device and
                    your account together.
                </p>
            </div>
        </Collapsable>

        <Collapsable
            text="Why do my favorites disappear when I login?"
            :button="false"
        >
            <div class="collapsable">
                <p>
                    When you're logged out, you're seeing the data that's saved
                    in your browser.
                </p>
                <p>
                    When you login, you're seeing the data that's saved in your
                    account.
                </p>
                <p>
                    If you want to see both, you'll need to sync them together.
                    Check the answer for "Can I sync up my device with my
                    account?"
                </p>
            </div>
        </Collapsable>

        <Collapsable text="What if I don't use an account?" :button="false">
            <div class="collapsable">
                <p>
                    Booruview will look and feel the exact same. However, your
                    data is only stored on your device. You can't share it
                    across other devices and it's at a greater risk of being
                    lost.
                </p>
            </div>
        </Collapsable>

        <!--
        <Collapsable text="Help, I forgot my password!" :button="false">
            <div class="collapsable">
                <p>
                    While you can't reset your password, I can try and help you
                    directly.
                </p>
                <p>
                    Send me an
                    <a href="mailto:feline8768@proton.me">email</a> and include:
                </p>
                <ul>
                    <li>your username</li>
                    <li>roughly how old your account is</li>
                    <li>when you last remember logging in</li>
                    <li>
                        a couple notable tags from your favs, blacklist, or
                        search history
                    </li>
                </ul>
            </div>
        </Collapsable>
        -->

        <Collapsable text="Will accounts ever be required?" :button="false">
            <div class="collapsable">
                <p>No, Booruview is free forever.</p>
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

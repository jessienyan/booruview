<script setup lang="ts">
import Collapsable from "@/components/Collapsable.vue";
</script>

<template>
    <div class="faq">
        <Collapsable text="What does having an account do?" :button="false">
            <div class="collapsable">
                <p>
                    Accounts lets you easily sync your data across your devices.
                    They're free, entirely optional, and no email required.
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

        <Collapsable text="I forgot my password." :button="false">
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

        <Collapsable text="Will accounts ever be required?" :button="false">
            <div class="collapsable">
                <p>No.</p>
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

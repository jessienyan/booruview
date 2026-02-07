<script setup lang="ts">
import store from "@/store";
</script>

<template>
	<div>
		<div v-if="store.account === null">
			<p>Accounts are free and totally optional.</p>

			<p>Without an account, your blacklist, favorites, and search history are all
				stored in your browser. You risk losing it if you wipe your browser data, and it's not shared
				between
				multiple devices.</p>
			<p>
			<form>
				<input class="input" type="text" />
				<input class="input" type="password" />
				<input class="input" type="password" />
			</form>
			</p>
		</div>
	</div>
</template>

<style lang="scss" scoped>
@import "@/assets/colors";
@import "@/assets/form";
</style>

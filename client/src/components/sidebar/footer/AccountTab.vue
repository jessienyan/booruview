<script setup lang="ts">
import { ref } from "vue";
import store from "@/store";

const username = ref<string | null>(null);
const password = ref<string | null>(null);
const passwordConfirm = ref<string | null>(null);
const formErr = ref("");

async function onSubmit(e: Event) {
	e.preventDefault();

	formErr.value = "";
	const $form = e.target as HTMLFormElement;

	if (!$form.reportValidity()) {
		return;
	} else if (password.value !== passwordConfirm.value) {
		formErr.value = "Password don't match";
		return;
	}

	const body = JSON.stringify({
		username: username.value,
		password: password.value,
	});

	try {
		const resp = await fetch("/api/register", {
			method: "POST",
			body,
			headers: {
				"Content-Type": "application/json",
			},
		});

		const data = await resp.json();
		if (data.error) {
			formErr.value = data.error;
			return;
		}
	} catch (e) {
		console.error(e);
	}
}
</script>

<template>
	<div>
		<div v-if="store.account === null">
			<p>Accounts are free and 100% optional.</p>
			<p>However, without an account, your blacklist, favorites, and search history are all stored in your browser. You risk losing it if you wipe your browser data.</p>

			<h2>register</h2>

			<form @submit="onSubmit">
				<input
					v-model="username"
					class="text-input rounded"
					:class="{touched: username != null}"
					type="text"
					placeholder="username"
					pattern="[a-zA-Z0-9_\-]+"
					maxlength="16"
					minlength="3"
					required />
				<span class="tip">Between 3 and 16 characters. Your username is private and only for logging in.</span>

				<input
					v-model="password"
					class="text-input rounded"
					:class="{touched: password != null}"
					type="password"
					placeholder="password"
					minlength="8"
					required />
				<span class="tip">At least 8 characters. YOU CANNOT RESET YOUR PASSWORD. Use something that's easy to remember!</span>

				<input
					v-model="passwordConfirm"
					class="text-input rounded"
					:class="{touched: passwordConfirm != null}"
					type="password"
					placeholder="confirm password"
					:minlength="8"
					required />

				<button class="submit btn-primary btn-rounded" type="submit">register</button>
			</form>
		</div>
	</div>
</template>

<style lang="scss" scoped>
@import "@/assets/buttons";
@import "@/assets/colors";
@import "@/assets/form";

form {
	padding: 0 1em;
}

.text-input {
	text-align: center;
	width: 100%;
}

.field-error, .tip {
	display: block;
	font-size: 14px;
	font-style: italic;
	margin: 0.5em 0.5em 1.5em;
}

.submit {
	margin-top: 1.5em;
	width: 100%;
	display: block;
}
</style>

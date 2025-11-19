import { useNewFeatureIndicator } from "@/composable";

export default {
	defaultBlacklist: useNewFeatureIndicator(
		"default-blacklist",
		new Date("2025-11-01"),
	),
};

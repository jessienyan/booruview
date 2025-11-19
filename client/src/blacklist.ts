export function defaultSFWBlacklist(): Tag[] {
	return [
		{ name: "rating:explicit", type: "unknown", count: 0 },
		{ name: "rating:questionable", type: "unknown", count: 0 },
		{ name: "rating:sensitive", type: "unknown", count: 0 },
	];
}

export function defaultNSFWBlacklist(): Tag[] {
	return [
		{ name: "all_the_way_through", type: "tag", count: 0 },
		{ name: "guro", type: "tag", count: 0 },
		{ name: "loli", type: "tag", count: 0 },
		{ name: "rape", type: "tag", count: 0 },
		{ name: "scat", type: "tag", count: 0 },
		{ name: "shota", type: "tag", count: 0 },
		{ name: "torture", type: "tag", count: 0 },
		{ name: "vomit", type: "tag", count: 0 },
	];
}

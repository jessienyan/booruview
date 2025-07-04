export type ChangeLogEntry = {
    date: string;
    changes: string[];
};

export const changeLog: ChangeLogEntry[] = [
    {
        date: "2025-07-04",
        changes: [
            "added this changelog",
            "fixed autocomplete not working when doing an OR search",
        ],
    },
    {
        date: "2025-07-01",
        changes: [
            'added the "recent" search tab',
            "fixed fullscreen view staying open when using back/forward browser nav",
        ],
    },
    {
        date: "2025-06-24",
        changes: ["added blacklist"],
    },
    {
        date: "2025-??-??",
        changes: [
            "too lazy to go back and write them all out",
            "check the commit history lol",
        ],
    },
    {
        date: "2025-05-11",
        changes: ["initial release"],
    },
];

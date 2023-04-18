pub const LISTS : [(&str, &str, &str); 16]= [
    // target
    ("list target", "list available targets in current workspace", "list target"),
    ("add target", "add a target to workspace", "add target name=NAME, type=TYPE"),
    ("del target", "delete a target", "del target id=ID"),
    ("set target", "update a target", "set target id=ID, name=NAME, type=TYPE, meta1=META1..."),
    // links
    ("list link", "list available links in current workspace", "list link"),
    ("add link", "add a link to workspace", "add link from=ID, to=ID, type=(IN/OUT/BOTH)"),
    ("del link", "delete a link", "del link id=ID"),
    ("set link", "update a link", "set link to=ID, from=NAME, meta1=META1..."),
    // modules
    ("list module", "list available module", "list module"),
    ("run module_name", "run a module", "run module_name, target_id=ID, arg1=..."),
    ("help module_name", "show help for specific module name", "help module_name"),
    // workspace
    ("list workspace", "list available workspace", "list workspace"),
    ("add workspace", "add a workspace", "add workspace name=NAME"),
    ("switch to workspace", "switch to another workspace", "switch workspace id=ID"),
    ("set workspace", "update a workspace name", "set workspace id=, name=NAME"),
    // export
    ("export dot", "export workspace in dot format", "export dot"),
];
pub fn extends_common(ctx: &mut rlua::Context) {
    let split_string = match ctx.create_function(|_, (separator, s): (String, String)| {
        let results = s.split(&separator);
        let results = results
            .collect::<Vec<&str>>()
            .iter()
            .map(|s| s.to_string())
            .collect::<Vec<String>>();
        Ok(results)
    }) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    let trim = match ctx.create_function(|_, s: String| Ok(s.trim().to_string())) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    ctx.globals().set("split", split_string).unwrap();
    ctx.globals().set("trim", trim).unwrap();
}

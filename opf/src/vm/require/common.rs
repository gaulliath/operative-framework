use std::collections::HashMap;

pub fn extends_common<'a>(ctx: &'a mut rlua::Context) {

    let create_target = match ctx.create_function(|_, (target): (HashMap<String, String>)| {
        Ok(true)
    }) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    ctx.globals().set("create_target", create_target).unwrap();
}

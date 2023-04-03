use crate::modules::metadata::Arg;
use rlua::{UserData, UserDataMethods};

pub struct Args {
    pub arg: Vec<Arg>,
}

impl From<Vec<Arg>> for Args {
    fn from(args: Vec<Arg>) -> Self {
        Self { arg: args.clone() }
    }
}

impl UserData for Args {
    fn add_methods<'lua, T: UserDataMethods<'lua, Self>>(methods: &mut T) {
        methods.add_method("get", |_, this, name: String| {
            for arg in &this.arg {
                if arg.name.eq(&name) {
                    return Ok(Some(arg.value.clone()));
                }
            }
            Ok(None)
        });
    }
}

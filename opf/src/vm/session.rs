use crate::core::session::Session;
use rlua::{UserData, UserDataMethods};

impl UserData for Session {

    fn add_methods<'lua, T: UserDataMethods<'lua, Self>>(methods: &mut T) {
        methods.add_method("get_id", |_, this, name: String| {
            Ok(this.session_id.clone().to_string())
        });
    }
}

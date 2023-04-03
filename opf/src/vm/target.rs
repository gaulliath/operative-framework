use crate::common::Target;
use rlua::{UserData, UserDataMethods};

impl UserData for Target {
    fn add_methods<'lua, T: UserDataMethods<'lua, Self>>(methods: &mut T) {
        methods.add_method("get_id", |_, this, ()| Ok(this.target_id.clone()));

        methods.add_method("get_name", |_, this, ()| Ok(this.target_name.clone()));

        methods.add_method("get_meta", |_, this, name: String| {
            if let Some(value) = this.meta.get(&name) {
                return Ok(Some(value.clone()));
            }
            Ok(None)
        });
    }
}

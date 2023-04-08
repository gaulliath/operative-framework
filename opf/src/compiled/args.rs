use crate::modules::metadata::Arg;

impl Arg {
    pub fn new(name: &str, is_target: bool, is_optional: bool, value: Option<String>) -> Self {
        Self {
            is_target,
            is_optional,
            name: name.to_string(),
            value,
        }
    }
}

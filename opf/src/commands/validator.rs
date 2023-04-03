use crate::error::ErrorKind;
use std::str::FromStr;
use strum_macros::EnumString;

#[derive(Debug, PartialEq, EnumString)]
#[strum(serialize_all = "lowercase")]
pub enum CommandAction {
    Add,
    Del,
    List,
    Set,
    Run,
    Accept,
    Save,
    Load,
    Export,
    Stop,
    Help,
}

#[derive(Debug, PartialEq, EnumString)]
#[strum(serialize_all = "lowercase")]
pub enum CommandObject {
    None,
    Save,
    Target,
    Link,
    Session,
    Action,
    Group,
    Api,
    Export(String),
    Module(String),
}

pub fn validate_action(action: String) -> Result<CommandAction, ErrorKind> {
    match CommandAction::from_str(action.to_lowercase().as_str()) {
        Ok(action) => Ok(action),
        Err(_) => Err(ErrorKind::ActionNotAvailable),
    }
}

pub fn validate_object(action: &CommandAction, object: String) -> Result<CommandObject, ErrorKind> {
    match CommandObject::from_str(object.to_lowercase().as_str()) {
        Ok(object) => Ok(object),
        Err(_) => {
            if object.len() > 0 {
                match action {
                    CommandAction::Run => return Ok(CommandObject::Module(object)),
                    CommandAction::Export => return Ok(CommandObject::Export(object)),
                    _ => {}
                }

                return Ok(CommandObject::Module(object));
            }
            Err(ErrorKind::ObjectNotAvailable)
        }
    }
}

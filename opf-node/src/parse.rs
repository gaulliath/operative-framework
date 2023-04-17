use std::collections::HashMap;
use std::str::FromStr;

use opf_models::error::ErrorKind;
use opf_models::{Command, CommandAction, CommandObject};

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

pub fn format(line: &str) -> Result<Command, ErrorKind> {
    let (action, object, params) = {
        let split = line.splitn(3, " ").collect::<Vec<&str>>();
        let split_len = split.len();

        if split_len < 1 {
            return Err(ErrorKind::InvalidCommandArgument);
        }

        let action = validate_action(split[0].to_string())?;

        let mut object = CommandObject::None;
        let mut params = HashMap::new();

        if split_len > 1 {
            object = validate_object(&action, split[1].to_string())?;

            if split_len > 2 {
                let parameters = split[2].split(",").collect::<Vec<&str>>();

                for param in parameters {
                    if !param.contains("=") {
                        return Err(ErrorKind::InvalidFormatArgument);
                    }

                    let vals = param.split("=").collect::<Vec<&str>>();
                    if vals.len() < 1 {
                        return Err(ErrorKind::InvalidFormatArgument);
                    }
                    params.insert(vals[0].trim().to_string(), vals[1].trim().to_string());
                }
            }
        }
        (action, object, params)
    };
    Ok(Command {
        action,
        object,
        params,
    })
}
